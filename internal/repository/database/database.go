package database

import (
	"log"
	"wb-tech-l0/internal/models"
	"wb-tech-l0/internal/repository/cache"

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
		&models.DeliveryDB{},
		&models.PaymentDB{},
		&models.OrderDB{},
		&models.ItemDB{},
	)
	return err
}

func (db *DB) SaveOrder(order *models.Order) error {
	err := db.Conn.Transaction(func(tx *gorm.DB) error {
		// Сохраняем Delivery
		deliveryDB := models.DeliveryDB{Delivery: order.Delivery}
		if err := tx.Create(&deliveryDB).Error; err != nil {
			return err
		}

		// Сохраняем Payment с ссылкой на OrderUID
		paymentDB := models.PaymentDB{
			Payment:  order.Payment,
			OrderUID: order.OrderUID,
		}
		if err := tx.Create(&paymentDB).Error; err != nil {
			return err
		}

		// Создаем OrderDB без вложенных структур
		orderDB := models.OrderDB{
			OrderUID:          order.OrderUID,
			TrackNumber:       order.TrackNumber,
			Entry:             order.Entry,
			Locale:            order.Locale,
			InternalSignature: order.InternalSignature,
			CustomerID:        order.CustomerID,
			DeliveryService:   order.DeliveryService,
			Shardkey:          order.Shardkey,
			SmID:              order.SmID,
			DateCreated:       order.DateCreated,
			OofShard:          order.OofShard,
			DeliveryID:        deliveryDB.ID,
			PaymentID:         paymentDB.ID,
		}

		if err := tx.Create(&orderDB).Error; err != nil {
			return err
		}

		// Сохраняем Items с ссылкой на OrderUID
		for _, item := range order.Items {
			itemDB := models.ItemDB{
				Item:     item,
				OrderUID: order.OrderUID,
			}
			if err := tx.Create(&itemDB).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Сохраняем заказ в кэш после успешного сохранения в БД
	db.Cache.Set(order.OrderUID, order)
	return nil
}

// GetOrder всегда работает через кэш
func (db *DB) GetOrder(orderUID string) (*models.Order, error) {
	// Сначала проверяем кэш
	if order, exists := db.Cache.Get(orderUID); exists {
		log.Printf("Order %s found in cache", orderUID)
		return order, nil
	}

	// Если нет в кэше, загружаем из БД
	order, err := db.loadOrderFromDB(orderUID)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш для будущих запросов
	db.Cache.Set(orderUID, order)
	return order, nil
}

// Приватный метод для загрузки заказа из БД
func (db *DB) loadOrderFromDB(orderUID string) (*models.Order, error) {
	var orderDB models.OrderDB
	if err := db.Conn.Where("order_uid = ?", orderUID).First(&orderDB).Error; err != nil {
		return nil, err
	}

	var deliveryDB models.DeliveryDB
	if err := db.Conn.First(&deliveryDB, orderDB.DeliveryID).Error; err != nil {
		return nil, err
	}

	var paymentDB models.PaymentDB
	if err := db.Conn.Where("order_uid = ?", orderUID).First(&paymentDB).Error; err != nil {
		return nil, err
	}

	var itemsDB []models.ItemDB
	if err := db.Conn.Where("order_uid = ?", orderUID).Find(&itemsDB).Error; err != nil {
		return nil, err
	}

	items := make([]models.Item, len(itemsDB))
	for i, itemDB := range itemsDB {
		items[i] = itemDB.Item
	}

	order := &models.Order{
		OrderUID:          orderDB.OrderUID,
		TrackNumber:       orderDB.TrackNumber,
		Entry:             orderDB.Entry,
		Delivery:          deliveryDB.Delivery,
		Payment:           paymentDB.Payment,
		Items:             items,
		Locale:            orderDB.Locale,
		InternalSignature: orderDB.InternalSignature,
		CustomerID:        orderDB.CustomerID,
		DeliveryService:   orderDB.DeliveryService,
		Shardkey:          orderDB.Shardkey,
		SmID:              orderDB.SmID,
		DateCreated:       orderDB.DateCreated,
		OofShard:          orderDB.OofShard,
	}

	return order, nil
}

func (db *DB) LoadAllOrders() error {
	var orderDBs []models.OrderDB
	if err := db.Conn.Find(&orderDBs).Error; err != nil {
		return err
	}

	for _, orderDB := range orderDBs {
		// Используем GetOrder, который работает через кэш
		order, err := db.GetOrder(orderDB.OrderUID)
		if err != nil {
			log.Printf("Failed to load order %s: %v", orderDB.OrderUID, err)
			continue
		}
		// GetOrder уже сохраняет заказ в кэш, поэтому ничего дополнительно делать не нужно
		_ = order // Используем переменную, чтобы избежать ошибки компиляции
	}

	log.Printf("Loaded %d orders to cache", db.Cache.Size())
	return nil
}

// GetOrderCount возвращает количество заказов в БД
func (db *DB) GetOrderCount() (int64, error) {
	var count int64
	if err := db.Conn.Model(&models.OrderDB{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
