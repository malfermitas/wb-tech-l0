package ports

import "wb-tech-l0/internal/models"

type OrderRepository interface {
	SaveOrder(order *models.Order) error
	GetOrder(orderUID string) (*models.Order, error)
	GetOrderCount() (int64, error)
}
