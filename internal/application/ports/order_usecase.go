package ports

import "wb-tech-l0/internal/models"

type OrderUseCase interface {
	ReceiveOrder(order *models.Order) error
	GetOrder(uid string) (*models.Order, error)
	GetStats() (OrderStats, error)
}

type OrderStats struct {
	CacheSize int
	DBCount   int64
}
