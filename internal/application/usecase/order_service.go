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

func (s *OrderService) Stats() (map[string]any, error) {
	count, err := s.repo.GetOrderCount()
	if err != nil {
		return nil, err
	}

	return map[string]any{"db_count": count}, nil
}
