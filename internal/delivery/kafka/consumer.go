package kafka

import (
	"encoding/json"
	"log"

	"wb-tech-l0/internal/application/ports"
	"wb-tech-l0/internal/models"

	"github.com/IBM/sarama"
)

// Consumer is a Kafka adapter that depends on the OrderUseCase, not on DB.
type Consumer struct {
	consumer     sarama.Consumer
	orderUseCase ports.OrderUseCase
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
	}, nil
}

func (c *Consumer) Start(topic string) error {
	partitionConsumer, err := c.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer partitionConsumer.Close()

	log.Println("Kafka consumer started. Waiting for messages...")

	for msg := range partitionConsumer.Messages() {
		log.Printf("Received message: %s\n", string(msg.Value))

		var order models.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("Error parsing JSON: %v\n", err)
			continue
		}

		if err := c.orderUseCase.ReceiveOrder(&order); err != nil {
			log.Printf("Failed to process order %s: %v\n", order.OrderUID, err)
		} else {
			log.Printf("Order %s processed successfully\n", order.OrderUID)
		}
	}

	return nil
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
