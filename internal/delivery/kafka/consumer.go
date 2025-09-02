package kafka

import (
	"encoding/json"
	"log"
	"time"
	"wb-tech-l0/internal/models"

	"github.com/IBM/sarama"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type OrderDB struct {
	gorm.Model
	OrderUID          string    `gorm:"primaryKey;uniqueIndex" json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`

	DeliveryID uint            `gorm:"not null" json:"-"`
	PaymentID  uint            `gorm:"not null" json:"-"`
	Delivery   models.Delivery `gorm:"-" json:"delivery"`
	Payment    models.Payment  `gorm:"-" json:"payment"`
	Items      []models.Item   `gorm:"-" json:"items"`
}

type DeliveryDB struct {
	gorm.Model
	models.Delivery
}

type PaymentDB struct {
	gorm.Model
	models.Payment
	OrderUID string `gorm:"not null;uniqueIndex" json:"-"`
}

type ItemDB struct {
	gorm.Model
	models.Item
	OrderUID string `gorm:"not null;index" json:"-"`
}

func Init_consumer() {
	dsn := "host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(&DeliveryDB{}, &PaymentDB{}, &OrderDB{}, &ItemDB{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatal("Failed to create consumer:", err)
	}
	defer consumer.Close()

	// Подписка на тему
	partitionConsumer, err := consumer.ConsumePartition("orders", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("Failed to create partition consumer:", err)
	}
	defer partitionConsumer.Close()

	log.Println("Consumer started. Waiting for messages...")

	for message := range partitionConsumer.Messages() {
		log.Printf("Received message: %s\n", string(message.Value))

		var order OrderDB
		if err := json.Unmarshal(message.Value, &order); err != nil {
			log.Printf("Error parsing JSON: %v\n", err)
			continue
		}

		// Сохранение в базу данных в транзакции
		err = db.Transaction(func(tx *gorm.DB) error {
			// Сохраняем Delivery
			deliveryDB := DeliveryDB{Delivery: order.Delivery}
			if err := tx.Create(&deliveryDB).Error; err != nil {
				return err
			}

			// Сохраняем Payment с ссылкой на OrderUID
			paymentDB := PaymentDB{
				Payment:  order.Payment,
				OrderUID: order.OrderUID,
			}
			if err := tx.Create(&paymentDB).Error; err != nil {
				return err
			}

			// Создаем OrderDB без вложенных структур
			orderDB := OrderDB{
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
				itemDB := ItemDB{
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
			log.Printf("Failed to save order %s: %v\n", order.OrderUID, err)
		} else {
			log.Printf("Order %s saved successfully\n", order.OrderUID)
		}
	}
}

func GetFullOrder(db *gorm.DB, orderUID string) (*OrderDB, error) {
	var orderDB OrderDB
	if err := db.Where("order_uid = ?", orderUID).First(&orderDB).Error; err != nil {
		return nil, err
	}

	var deliveryDB DeliveryDB
	if err := db.First(&deliveryDB, orderDB.DeliveryID).Error; err != nil {
		return nil, err
	}

	var paymentDB PaymentDB
	if err := db.Where("order_uid = ?", orderUID).First(&paymentDB).Error; err != nil {
		return nil, err
	}

	var itemsDB []ItemDB
	if err := db.Where("order_uid = ?", orderUID).Find(&itemsDB).Error; err != nil {
		return nil, err
	}

	items := make([]models.Item, len(itemsDB))
	for i, itemDB := range itemsDB {
		items[i] = itemDB.Item
	}

	order := &OrderDB{
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
