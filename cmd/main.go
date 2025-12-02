package main

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"wb-tech-l0/internal/config"

	"wb-tech-l0/cmd/server"
	"wb-tech-l0/internal/application/usecase"
	"wb-tech-l0/internal/delivery/kafka"
	"wb-tech-l0/internal/repository/cache"
	"wb-tech-l0/internal/repository/database"

	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	// Root context with OS signal cancellation.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// --- Infrastructure setup ---

	redisClient := newRedisClient(cfg.RedisAddr)
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("Failed to close Redis client: %v", err)
		}
	}()

	orderCache := cache.NewOrderCache(redisClient, cfg.CacheTTL)

	db := newDatabase(cfg.PostgresDSN, orderCache)
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
		cfg.KafkaBrokers,
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
	wg.Add(3)

	// HTTP server lifecycle
	go func() {
		defer wg.Done()
		log.Printf("Starting HTTP server on %s", cfg.HTTPAddr)
		if err := httpServer.Start(cfg.HTTPAddr); err != nil {
			// http.ErrServerClosed is expected on graceful shutdown
			log.Printf("HTTP server stopped: %v", err)
		}
	}()

	// Kafka consumer lifecycle
	go func() {
		defer wg.Done()
		log.Printf("Starting Kafka consumer on %s, topic %s", cfg.KafkaBrokers, cfg.KafkaTopic)
		if err := kafkaConsumer.Start(ctx, cfg.KafkaTopic); err != nil && ctx.Err() == nil {
			log.Printf("Kafka consumer error: %v", err)
		}
	}()

	// Fill cache with orders from DB

	go func() {
		defer wg.Done()
		log.Printf("Filling cache with orders from DB...")
		if err := orderUC.LoadOrdersToCache(cfg.CachePreloadCount); err != nil {
			log.Printf("Failed to load orders to cache: %v", err)
		}
	}()

	// Graceful shutdown implementation

	<-ctx.Done()
	log.Println("Shutdown signal received, shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	wg.Wait()

	log.Println("Shutdown complete")
}

func newRedisClient(addr string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis at %s: %v", addr, err)
	}

	log.Printf("Connected to Redis at %s", addr)

	return client
}

func newDatabase(dsn string, orderCache *cache.OrderCache) *database.DB {
	db, err := database.NewDB(dsn, orderCache)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Connected to database")

	return db
}
