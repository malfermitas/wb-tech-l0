package ports

import "wb-tech-l0/internal/models"

type OrderUseCase interface {
	SaveOrder(order *models.Order) error
	GetOrder(uid string) (*models.Order, error)
	Stats() (OrderStats, error)
}

type OrderStats struct {
	CacheSize int
	DBCount   int64
}
