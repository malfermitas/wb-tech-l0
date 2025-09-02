package cache

import (
	"log"
	"sync"
	"wb-tech-l0/internal/models"
)

type OrderCache struct {
	sync.RWMutex
	orders map[string]*models.Order
}

func NewOrderCache() *OrderCache {
	return &OrderCache{
		orders: make(map[string]*models.Order),
	}
}

func (c *OrderCache) Set(orderUID string, order *models.Order) {
	c.Lock()
	defer c.Unlock()
	c.orders[orderUID] = order
	log.Println("Order", orderUID, "has been set to cache")
}

func (c *OrderCache) Get(orderUID string) (*models.Order, bool) {
	c.RLock()
	defer c.RUnlock()
	order, exists := c.orders[orderUID]
	return order, exists
}

func (c *OrderCache) GetAll() map[string]*models.Order {
	c.RLock()
	defer c.RUnlock()

	// Создаем копию для безопасного возврата
	result := make(map[string]*models.Order)
	for k, v := range c.orders {
		result[k] = v
	}
	return result
}

func (c *OrderCache) Preload(orders map[string]*models.Order) {
	c.Lock()
	defer c.Unlock()
	for k, v := range orders {
		c.orders[k] = v
	}
}

func (c *OrderCache) Remove(orderUID string) {
	c.Lock()
	defer c.Unlock()
	delete(c.orders, orderUID)
}

func (c *OrderCache) Clear() {
	c.Lock()
	defer c.Unlock()
	c.orders = make(map[string]*models.Order)
}

func (c *OrderCache) Size() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.orders)
}
