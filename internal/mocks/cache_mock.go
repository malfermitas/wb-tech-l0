package mocks

import (
	"time"
	"wb-tech-l0/internal/models"

	"github.com/stretchr/testify/mock"
)

type CacheMock struct {
	mock.Mock
}

func (m *CacheMock) Set(key string, order *models.Order, ttl time.Duration) error {
	args := m.Called(key, order, ttl)
	return args.Error(0)
}

func (m *CacheMock) Get(key string) (*models.Order, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *CacheMock) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}
