package kafka

import (
	"context"
	"encoding/json"
	"log"
	"wb-tech-l0/internal/validator"

	"wb-tech-l0/internal/application/ports"
	"wb-tech-l0/internal/models"

	"github.com/IBM/sarama"
)

// Consumer is a Kafka adapter that depends on the OrderUseCase, not on DB.
type Consumer struct {
	consumer     sarama.Consumer
	orderUseCase ports.OrderUseCase
	validator    validator.Validator
}

func NewConsumer(brokers []string, uc ports.OrderUseCase) (*Consumer, error) {
	cfg := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(brokers, cfg)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer:     consumer,
		orderUseCase: uc,
		validator:    validator.NewValidator(),
	}, nil
}

// NewConsumerWith allows injecting a custom sarama.Consumer and validator, making it test-friendly.
func NewConsumerWith(consumer sarama.Consumer, uc ports.OrderUseCase, v validator.Validator) *Consumer {
	return &Consumer{
		consumer:     consumer,
		orderUseCase: uc,
		validator:    v,
	}
}

// Start consumes messages from the given topic until the context is cancelled.
func (c *Consumer) Start(ctx context.Context, topic string) error {
	partitionConsumer, err := c.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer partitionConsumer.Close()

	log.Println("Kafka consumer started. Waiting for messages...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer: context cancelled, stopping")
			return nil

		case msg, ok := <-partitionConsumer.Messages():
			if !ok {
				log.Println("Kafka consumer: partition consumer channel closed")
				return nil
			}

			log.Printf("Received message: %s\n", string(msg.Value))

			// Парсинг JSON
			var order models.Order
			if err := json.Unmarshal(msg.Value, &order); err != nil {
				log.Printf("Error parsing JSON: %v\n", err)
				continue
			}

			// Валидация модели

			if err := c.validator.Validate(order); err != nil {
				log.Printf("❌ Invalid order data for orderUID %s: %v", order.OrderUID, err)
				return err
			}

			if err := c.orderUseCase.SaveOrder(&order); err != nil {
				log.Printf("Failed to process order %s: %v\n", order.OrderUID, err)
				return err
			} else {
				log.Printf("Order %s processed successfully\n", order.OrderUID)
			}
		}
	}
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
