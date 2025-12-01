package main

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"wb-tech-l0/cmd/server"
	"wb-tech-l0/internal/application/usecase"
	"wb-tech-l0/internal/delivery/kafka"
	"wb-tech-l0/internal/repository/cache"
	"wb-tech-l0/internal/repository/database"

	"github.com/redis/go-redis/v9"
)

const (
	httpAddr      = ":8080"
	kafkaTopic    = "orders"
	kafkaBroker   = "localhost:9092"
	postgresDSN   = "host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable"
	redisAddr     = "localhost:6379"
	cacheTTL      = 10 * time.Minute
	shutdownGrace = 10 * time.Second
)

func main() {
	// Root context with OS signal cancellation.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// --- Infrastructure setup ---

	redisClient := newRedisClient()
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("Failed to close Redis client: %v", err)
		}
	}()

	orderCache := cache.NewOrderCache(redisClient, cacheTTL)

	db := newDatabase(orderCache)
	defer func() {
		sqlDB, err := db.Conn.DB()
		if err != nil {
			log.Printf("Failed to get sql.DB from GORM: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("Failed to close DB connection: %v", err)
		}
	}()

	// --- Application layer ---

	orderRepo := db
	orderUC := usecase.NewOrderService(orderRepo)

	// --- Delivery / adapters ---

	httpServer := server.NewServer(orderUC)

	kafkaConsumer, err := kafka.NewConsumer(
		[]string{kafkaBroker},
		orderUC,
	)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer func() {
		if err := kafkaConsumer.Close(); err != nil {
			log.Printf("Failed to close Kafka consumer: %v", err)
		}
	}()

	// --- Run servers ---

	var wg sync.WaitGroup
	wg.Add(2)

	// HTTP server lifecycle
	go func() {
		defer wg.Done()
		log.Printf("Starting HTTP server on %s", httpAddr)
		if err := httpServer.Start(httpAddr); err != nil {
			// http.ErrServerClosed is expected on graceful shutdown
			log.Printf("HTTP server stopped: %v", err)
		}
	}()

	// Kafka consumer lifecycle
	go func() {
		defer wg.Done()
		log.Printf("Starting Kafka consumer on %s, topic %s", kafkaBroker, kafkaTopic)
		if err := kafkaConsumer.Start(ctx, kafkaTopic); err != nil && ctx.Err() == nil {
			log.Printf("Kafka consumer error: %v", err)
		}
	}()

	// Graceful shutdown implementation

	<-ctx.Done()
	log.Println("Shutdown signal received, shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownGrace)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	wg.Wait()

	log.Println("Shutdown complete")
}

func newRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		// Add Password / DB from config if needed.
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis at %s: %v", redisAddr, err)
	}

	return client
}

func newDatabase(orderCache *cache.OrderCache) *database.DB {
	db, err := database.NewDB(postgresDSN, orderCache)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}
