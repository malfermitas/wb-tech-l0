package mocks

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/mock"
)

// PartitionConsumerMock реализует sarama.PartitionConsumer.
type PartitionConsumerMock struct {
	mock.Mock
	msgs chan *sarama.ConsumerMessage
	errs chan *sarama.ConsumerError
}

func (m *PartitionConsumerMock) Pause() {
	m.Called()
}

func (m *PartitionConsumerMock) Resume() {
	m.Called()
}

func (m *PartitionConsumerMock) IsPaused() bool {
	m.Called()
	return false
}

var _ sarama.PartitionConsumer = (*PartitionConsumerMock)(nil)

func NewPartitionConsumerMock() *PartitionConsumerMock {
	return &PartitionConsumerMock{
		msgs: make(chan *sarama.ConsumerMessage, 100),
		errs: make(chan *sarama.ConsumerError, 10),
	}
}

func (m *PartitionConsumerMock) AsyncClose() {
	m.Called()
	close(m.msgs)
	close(m.errs)
}

func (m *PartitionConsumerMock) Close() error {
	args := m.Called()
	close(m.msgs)
	close(m.errs)
	return args.Error(0)
}

func (m *PartitionConsumerMock) Messages() <-chan *sarama.ConsumerMessage {
	return m.msgs
}

func (m *PartitionConsumerMock) Errors() <-chan *sarama.ConsumerError {
	return m.errs
}

func (m *PartitionConsumerMock) HighWaterMarkOffset() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

// Вспомогательные методы для тестов.
func (m *PartitionConsumerMock) PushMessage(msg *sarama.ConsumerMessage) {
	m.msgs <- msg
}

func (m *PartitionConsumerMock) PushError(err *sarama.ConsumerError) {
	m.errs <- err
}
