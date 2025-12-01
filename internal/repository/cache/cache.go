package cache

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"wb-tech-l0/internal/models"

	"github.com/redis/go-redis/v9"
)

type OrderCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewOrderCache(client *redis.Client, ttl time.Duration) *OrderCache {
	return &OrderCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *OrderCache) key(orderUID string) string {
	return "order:" + orderUID
}

func (c *OrderCache) Set(orderUID string, order *models.Order) {
	ctx := context.Background()

	data, err := json.Marshal(order)
	if err != nil {
		log.Printf("OrderCache: failed to marshal order %s: %v", orderUID, err)
		return
	}

	if err := c.client.Set(ctx, c.key(orderUID), data, c.ttl).Err(); err != nil {
		log.Printf("OrderCache: failed to set order %s in Redis: %v", orderUID, err)
	}
}

func (c *OrderCache) Get(orderUID string) (*models.Order, bool) {
	ctx := context.Background()

	val, err := c.client.Get(ctx, c.key(orderUID)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, false
	}
	if err != nil {
		log.Printf("OrderCache: failed to get order %s from Redis: %v", orderUID, err)
		return nil, false
	}

	var order models.Order
	if err := json.Unmarshal([]byte(val), &order); err != nil {
		log.Printf("OrderCache: failed to unmarshal order %s from Redis: %v", orderUID, err)
		return nil, false
	}

	return &order, true
}

func (c *OrderCache) Size() int {
	ctx := context.Background()

	n, err := c.client.DBSize(ctx).Result()
	if err != nil {
		log.Printf("OrderCache: failed to get DB size from Redis: %v", err)
		return 0
	}

	return int(n)
}
