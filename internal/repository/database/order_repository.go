package database

import (
	"log"
	"wb-tech-l0/internal/application/ports"
	"wb-tech-l0/internal/models"
	"wb-tech-l0/internal/repository/database/db_models"

	"gorm.io/gorm"
)

var _ ports.OrderRepository = (*DB)(nil)

func (db *DB) SaveOrder(order *models.Order) error {
	err := db.Conn.Transaction(func(tx *gorm.DB) error {
		deliveryDB := db_models.ToDeliveryDB(order.Delivery)
		if err := tx.Create(&deliveryDB).Error; err != nil {
			return err
		}

		paymentDB := db_models.ToPaymentDB(order)
		if err := tx.Create(&paymentDB).Error; err != nil {
			return err
		}

		orderDB := db_models.ToOrderDB(order, deliveryDB.ID, paymentDB.ID)
		if err := tx.Create(&orderDB).Error; err != nil {
			return err
		}

		for _, item := range order.Items {
			itemDB := db_models.ToItemDB(item, order.OrderUID)
			if err := tx.Create(&itemDB).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	db.Cache.Set(order.OrderUID, order)
	return nil
}

func (db *DB) GetOrder(orderUID string) (*models.Order, error) {
	if order, ok := db.Cache.Get(orderUID); ok {
		log.Printf("Order %s found in cache", orderUID)
		return order, nil
	}

	order, err := db.loadOrderFromDB(orderUID)
	if err != nil {
		return nil, err
	}

	db.Cache.Set(orderUID, order)
	return order, nil
}

func (db *DB) loadOrderFromDB(orderUID string) (*models.Order, error) {
	var orderDB db_models.OrderDB
	if err := db.Conn.Where("order_uid = ?", orderUID).First(&orderDB).Error; err != nil {
		return nil, err
	}

	var deliveryDB db_models.DeliveryDB
	if err := db.Conn.First(&deliveryDB, orderDB.DeliveryID).Error; err != nil {
		return nil, err
	}

	var paymentDB db_models.PaymentDB
	if err := db.Conn.Where("order_uid = ?", orderUID).First(&paymentDB).Error; err != nil {
		return nil, err
	}

	var itemsDB []db_models.ItemDB
	if err := db.Conn.Where("order_uid = ?", orderUID).Find(&itemsDB).Error; err != nil {
		return nil, err
	}

	return db_models.ToDomainOrder(orderDB, deliveryDB, paymentDB, itemsDB), nil
}

func (db *DB) LoadAllOrdersToCache() error {
	var orderDBs []db_models.OrderDB
	if err := db.Conn.Find(&orderDBs).Error; err != nil {
		return err
	}

	for _, odb := range orderDBs {
		order, err := db.GetOrder(odb.OrderUID)
		if err != nil {
			log.Printf("Failed to load order %s: %v", odb.OrderUID, err)
			continue
		}
		_ = order
	}

	log.Printf("Loaded %d orders to cache", db.Cache.Size())
	return nil
}

func (db *DB) GetOrderCount() (int64, error) {
	var count int64
	if err := db.Conn.Model(&db_models.OrderDB{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (db *DB) CacheSize() int {
	return db.Cache.Size()
}
