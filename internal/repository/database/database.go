package database

import (
	"wb-tech-l0/internal/repository/cache"
	"wb-tech-l0/internal/repository/database/db_models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	Conn  *gorm.DB
	Cache *cache.OrderCache
}

func NewDB(dsn string, cache *cache.OrderCache) (*DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &DB{
		Conn:  db,
		Cache: cache,
	}, nil
}

func (db *DB) Migrate() error {
	err := db.Conn.AutoMigrate(
		&db_models.DeliveryDB{},
		&db_models.PaymentDB{},
		&db_models.OrderDB{},
		&db_models.ItemDB{},
	)
	return err
}
