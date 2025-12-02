package mocks

import (
	"wb-tech-l0/internal/application/ports"
	"wb-tech-l0/internal/models"

	"github.com/stretchr/testify/mock"
)

// OrderRepositoryMock реализует интерфейс ports.OrderRepository.
type OrderRepositoryMock struct {
	mock.Mock
}

var _ ports.OrderRepository = (*OrderRepositoryMock)(nil)

func (m *OrderRepositoryMock) SaveOrder(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *OrderRepositoryMock) GetOrder(orderUID string) (*models.Order, error) {
	args := m.Called(orderUID)
	if v := args.Get(0); v != nil {
		return v.(*models.Order), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *OrderRepositoryMock) GetOrderCount() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *OrderRepositoryMock) CacheSize() int {
	args := m.Called()
	return args.Int(0)
}
