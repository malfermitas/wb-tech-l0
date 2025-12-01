package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"wb-tech-l0/internal/models"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

func main() {
	v := viper.New()
	v.AutomaticEnv()

	ordersCount := v.GetInt("ORDERS_COUNT")
	if ordersCount <= 0 {
		log.Fatalf("environment variable ORDERS_COUNT is required and must be > 0 (got %d)", ordersCount)
	}

	log.Printf("Генерируем %d заказов...\n", ordersCount)

	kafkaHost := v.GetString("KAFKA_HOST")
	if kafkaHost == "" {
		panic("environment variable KAFKA_HOST is required but not set")
	}

	topic := v.GetString("KAFKA_TOPIC")
	if topic == "" {
		panic("environment variable KAFKA_TOPIC is required but not set")
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaHost),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	var successCount int64
	var errorCount int64

	var wg sync.WaitGroup

	kafkaChan := make(chan int)

	numWorkers := 100

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for orderNum := range kafkaChan {
				order := generateOrder()
				err := sendToKafka(writer, order)

				if err != nil {
					log.Printf("❌ Worker %d: Ошибка отправки заказа %d: %v", workerID, orderNum, err)
					atomic.AddInt64(&errorCount, 1)
				} else {
					log.Printf("✅ Worker %d: Заказ %d отправлен: %s", workerID, orderNum, order.OrderUID)
					atomic.AddInt64(&successCount, 1)
				}
			}
		}(i)
	}

	for i := 0; i < ordersCount; i++ {
		kafkaChan <- i
	}
	close(kafkaChan)
	wg.Wait()

	fmt.Printf("\n=== Результаты ===\n")
	fmt.Printf("Успешно отправлено: %d\n", atomic.LoadInt64(&successCount))
	fmt.Printf("Ошибок: %d\n", atomic.LoadInt64(&errorCount))
	fmt.Printf("Тема Kafka: %s\n", topic)
	fmt.Printf("Брокер: %s\n", kafkaHost)
}

func generateOrder() models.Order {
	var order models.Order
	err := gofakeit.Struct(&order)
	if err != nil {
		log.Fatalf("Ошибка генерации заказа: %v", err)
	}
	order.Payment.Transaction = order.OrderUID
	for i := range order.Items {
		order.Items[i].TrackNumber = order.TrackNumber
	}

	return order
}

func sendToKafka(writer *kafka.Writer, order models.Order) error {
	jsonData, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("ошибка маршалинга JSON: %v", err)
	}

	err = writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(order.OrderUID),
			Value: jsonData,
		},
	)

	if err != nil {
		return fmt.Errorf("ошибка отправки в Kafka: %v", err)
	}

	return nil
}
