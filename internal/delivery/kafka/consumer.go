package kafka

import (
	"encoding/json"
	"log"
	"wb-tech-l0/internal/models"
	"wb-tech-l0/internal/repository/database"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumer sarama.Consumer
	db       *database.DB
}

func NewConsumer(brokers []string, db *database.DB) (*Consumer, error) {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: consumer,
		db:       db,
	}, nil
}

func (c *Consumer) Start(topic string) error {
	partitionConsumer, err := c.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer partitionConsumer.Close()

	log.Println("Kafka consumer started. Waiting for messages...")

	for message := range partitionConsumer.Messages() {
		log.Printf("Received message: %s\n", string(message.Value))

		var order models.Order
		if err := json.Unmarshal(message.Value, &order); err != nil {
			log.Printf("Error parsing JSON: %v\n", err)
			continue
		}

		// Сохранение в базу данных (автоматически сохраняет в кэш)
		if err := c.db.SaveOrder(&order); err != nil {
			log.Printf("Failed to save order %s: %v\n", order.OrderUID, err)
		} else {
			log.Printf("Order %s saved successfully\n", order.OrderUID)
		}
	}

	return nil
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
