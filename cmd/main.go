package main

import (
	"log"
	"time"
	"wb-tech-l0/cmd/server"
	"wb-tech-l0/internal/delivery/kafka"
	"wb-tech-l0/internal/repository/cache"
	"wb-tech-l0/internal/repository/database"
)

func main() {
	go func() {
		time.Sleep(3 * time.Second)
		server.RunFakeDataProducer()
	}()

	cache := cache.NewOrderCache()

	dsn := "host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable"
	db, err := database.NewDB(dsn, cache)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.Migrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	server := server.NewServer(db)
	go func() {
		log.Println("Starting HTTP server on :8080")
		if err := server.Start(":8080"); err != nil {
			log.Fatal("HTTP server error:", err)
		}
	}()

	kafkaConsumer, err := kafka.NewConsumer(
		[]string{"localhost:9092"},
		db,
	)
	if err != nil {
		log.Fatal("Failed to create Kafka consumer:", err)
	}
	defer kafkaConsumer.Close()

	if err := kafkaConsumer.Start("orders"); err != nil {
		log.Fatal("Kafka consumer error:", err)
	}
}
