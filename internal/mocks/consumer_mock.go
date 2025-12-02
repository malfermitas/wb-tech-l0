package mocks

import (
	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
	"github.com/stretchr/testify/mock"
)

// ConsumerMock реализует sarama.Consumer.
type ConsumerMock struct {
	mock.Mock
	*mocks.Consumer
}

// Ensure ConsumerMock satisfies sarama.Consumer at compile‑time.
var _ sarama.Consumer = (*ConsumerMock)(nil)

func NewConsumerMock() *ConsumerMock {
	return &ConsumerMock{
		Consumer: mocks.NewConsumer(nil, nil),
	}
}

func (c *ConsumerMock) ExpectConsumePartition(topic string, partition int32, offset int64) *mocks.PartitionConsumer {
	return c.Consumer.ExpectConsumePartition(topic, partition, offset)
}
