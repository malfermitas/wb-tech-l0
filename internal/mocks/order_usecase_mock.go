package mocks

import (
	"wb-tech-l0/internal/application/ports"
	"wb-tech-l0/internal/models"

	"github.com/stretchr/testify/mock"
)

// OrderUseCaseMock реализует интерфейс ports.OrderUseCase.
type OrderUseCaseMock struct {
	mock.Mock
}

// Компилятор проверит, что структура реализует интерфейс.
var _ ports.OrderUseCase = (*OrderUseCaseMock)(nil)

func (m *OrderUseCaseMock) SaveOrder(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *OrderUseCaseMock) GetOrder(orderUID string) (*models.Order, error) {
	args := m.Called(orderUID)
	if v := args.Get(0); v != nil {
		return v.(*models.Order), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *OrderUseCaseMock) Stats() (ports.OrderStats, error) {
	args := m.Called()
	var stats ports.OrderStats
	if v := args.Get(0); v != nil {
		stats = v.(ports.OrderStats)
	}
	return stats, args.Error(1)
}
