package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"wb-tech-l0/internal/models"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/segmentio/kafka-go"
)

type Order models.Order

func run_fake_data_producer() {
	kafkaURL := "localhost:9092"
	topic := "orders"

	// Создаем писателя в Kafka
	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	var wg sync.WaitGroup
	for i := 0; i < 50000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			order := generateOrder()
			err := sendToKafka(writer, order)
			if err != nil {
				log.Printf("Error sending order %d: %v", i, err)
			} else {
				log.Printf("Order %d has been sent", i)
			}
		}()
	}
	wg.Wait()
}

func generateOrder() Order {
	var order Order
	err := gofakeit.Struct(&order)
	if err != nil {
		log.Fatalf("Error generating order: %v", err)
	}

	order.Payment.Transaction = order.OrderUID
	for i := range order.Items {
		order.Items[i].TrackNumber = order.TrackNumber
	}

	return order
}

func sendToKafka(writer *kafka.Writer, order Order) error {
	jsonData, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	err = writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(order.OrderUID),
			Value: jsonData,
		},
	)

	if err != nil {
		return fmt.Errorf("error sending to Kafka: %v", err)
	}

	fmt.Printf("Order sent: %s\n", order.OrderUID)
	return nil
}
