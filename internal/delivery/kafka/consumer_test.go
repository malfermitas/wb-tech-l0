package kafka

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"wb-tech-l0/internal/application/ports"
	imocks "wb-tech-l0/internal/mocks"
	"wb-tech-l0/internal/models"
	"wb-tech-l0/internal/validator"

	"github.com/IBM/sarama"
	smocks "github.com/IBM/sarama/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// helper to start consumer with injected deps
func newTestConsumer(saramaConsumer sarama.Consumer, uc ports.OrderUseCase, v validator.Validator) *Consumer {
	return NewConsumerWith(saramaConsumer, uc, v)
}

func TestConsumer_SuccessfulProcessing(t *testing.T) {
	topic := "orders"
	saramaC := smocks.NewConsumer(nil, nil)
	defer saramaC.Close()

	pc := saramaC.ExpectConsumePartition(topic, 0, sarama.OffsetNewest)

	// Mocks for use case and validator
	uc := new(imocks.OrderUseCaseMock)
	v := new(imocks.ValidatorMock)

	// Arrange order and expectations
	order := models.Order{
		OrderUID:        "uid-1",
		TrackNumber:     "ABCDEFGHJK",
		Entry:           "WBIL",
		Locale:          "en",
		CustomerID:      "11111111-1111-1111-1111-111111111111",
		DeliveryService: "meest",
		Shardkey:        "1",
		SmID:            1,
		DateCreated:     time.Now().Truncate(time.Second),
		OofShard:        "1",
		Delivery:        models.Delivery{Name: "John", Phone: "+123", Zip: "12345", City: "City", Address: "Street 1", Region: "Reg", Email: "john@example.com"},
		Payment:         models.Payment{Transaction: "tx1", Currency: "USD", Provider: "wbpay", Amount: 100, PaymentDt: time.Now().Unix(), Bank: "bank", DeliveryCost: 10, GoodsTotal: 90, CustomFee: 0},
		Items:           []models.Item{{ChrtID: 1, TrackNumber: "ABCDEFGHJK", Price: 100, RID: "rid", Name: "Item", Sale: 0, Size: "M", TotalPrice: 100, NmID: 1, Brand: "brand", Status: 200}},
	}

	data, _ := json.Marshal(order)

	v.On("Validate", mock.Anything).Return(nil)
	uc.On("SaveOrder", mock.MatchedBy(func(o *models.Order) bool { return o != nil && o.OrderUID == order.OrderUID })).Return(nil)

	cons := newTestConsumer(saramaC, uc, v)

	ctx, cancel := context.WithCancel(context.Background())
	doneCh := make(chan error, 1)
	go func() { doneCh <- cons.Start(ctx, topic) }()

	// Send message then stop
	pc.YieldMessage(&sarama.ConsumerMessage{Value: data})
	time.Sleep(50 * time.Millisecond)
	cancel()

	err := <-doneCh
	assert.NoError(t, err)
	uc.AssertCalled(t, "SaveOrder", mock.MatchedBy(func(o *models.Order) bool { return o != nil && o.OrderUID == order.OrderUID }))
}

func TestConsumer_InvalidJSON(t *testing.T) {
	topic := "orders"
	saramaC := smocks.NewConsumer(nil, nil)
	defer saramaC.Close()
	pc := saramaC.ExpectConsumePartition(topic, 0, sarama.OffsetNewest)

	uc := new(imocks.OrderUseCaseMock)
	v := new(imocks.ValidatorMock)

	cons := newTestConsumer(saramaC, uc, v)
	ctx, cancel := context.WithCancel(context.Background())
	doneCh := make(chan error, 1)
	go func() { doneCh <- cons.Start(ctx, topic) }()

	pc.YieldMessage(&sarama.ConsumerMessage{Value: []byte("not-json")})
	time.Sleep(50 * time.Millisecond)
	cancel()
	err := <-doneCh
	assert.NoError(t, err)

	// Ensure validator and usecase were not called
	v.AssertNotCalled(t, "Validate", mock.Anything)
	uc.AssertNotCalled(t, "SaveOrder", mock.Anything)
}

func TestConsumer_ValidatorError(t *testing.T) {
	topic := "orders"
	saramaC := smocks.NewConsumer(nil, nil)
	defer saramaC.Close()
	pc := saramaC.ExpectConsumePartition(topic, 0, sarama.OffsetNewest)

	uc := new(imocks.OrderUseCaseMock)
	v := new(imocks.ValidatorMock)
	cons := newTestConsumer(saramaC, uc, v)

	order := models.Order{OrderUID: "uid-2"}
	data, _ := json.Marshal(order)

	v.On("Validate", mock.Anything).Return(assert.AnError)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errCh := make(chan error, 1)
	go func() { errCh <- cons.Start(ctx, topic) }()

	pc.YieldMessage(&sarama.ConsumerMessage{Value: data})
	err := <-errCh
	assert.Error(t, err)
	uc.AssertNotCalled(t, "SaveOrder", mock.Anything)
}

func TestConsumer_ContextCancel(t *testing.T) {
	topic := "orders"
	saramaC := smocks.NewConsumer(nil, nil)
	defer saramaC.Close()
	saramaC.ExpectConsumePartition(topic, 0, sarama.OffsetNewest)

	uc := new(imocks.OrderUseCaseMock)
	v := new(imocks.ValidatorMock)
	cons := newTestConsumer(saramaC, uc, v)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- cons.Start(ctx, topic) }()

	// cancel immediately
	cancel()
	err := <-done
	assert.NoError(t, err)
}
