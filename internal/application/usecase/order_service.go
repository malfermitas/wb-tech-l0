package usecase

import (
	"wb-tech-l0/internal/application/ports"
	"wb-tech-l0/internal/models"
)

type OrderService struct {
	repo ports.OrderRepository
}

func NewOrderService(repo ports.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) GetOrder(uid string) (*models.Order, error) {
	return s.repo.GetOrder(uid)
}

func (s *OrderService) SaveOrder(order *models.Order) error {
	return s.repo.SaveOrder(order)
}

func (s *OrderService) Stats() (ports.OrderStats, error) {
	dbCount, err := s.repo.GetOrderCount()
	if err != nil {
		return ports.OrderStats{}, err
	}

	return ports.OrderStats{
		CacheSize: s.repo.CacheSize(),
		DBCount:   dbCount,
	}, nil
}
