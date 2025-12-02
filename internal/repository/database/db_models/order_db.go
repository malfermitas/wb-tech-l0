package db_models

import (
	"time"
	"wb-tech-l0/internal/models"

	"gorm.io/gorm"
)

type OrderDB struct {
	gorm.Model
	OrderUID          string `gorm:"primaryKey;uniqueIndex"`
	TrackNumber       string
	Entry             string
	Locale            string
	InternalSignature string
	CustomerID        string
	DeliveryService   string
	Shardkey          string
	SmID              int
	DateCreated       int64
	OofShard          string

	DeliveryID uint `gorm:"not null"`
	PaymentID  uint `gorm:"not null"`
}

func ToOrderDB(o *models.Order, deliveryID, paymentID uint) OrderDB {
	return OrderDB{
		OrderUID:          o.OrderUID,
		TrackNumber:       o.TrackNumber,
		Entry:             o.Entry,
		Locale:            o.Locale,
		InternalSignature: o.InternalSignature,
		CustomerID:        o.CustomerID,
		DeliveryService:   o.DeliveryService,
		Shardkey:          o.Shardkey,
		SmID:              o.SmID,
		DateCreated:       o.DateCreated.Unix(),
		OofShard:          o.OofShard,
		DeliveryID:        deliveryID,
		PaymentID:         paymentID,
	}
}

func ToDomainOrder(
	orderDB OrderDB,
	deliveryDB DeliveryDB,
	paymentDB PaymentDB,
	itemsDB []ItemDB,
) *models.Order {
	items := make([]models.Item, len(itemsDB))
	for i, it := range itemsDB {
		items[i] = models.Item{
			ChrtID:      it.ChrtID,
			TrackNumber: it.TrackNumber,
			Price:       it.Price,
			RID:         it.Rid,
			Name:        it.Name,
			Sale:        it.Sale,
			Size:        it.Size,
			TotalPrice:  it.TotalPrice,
			NmID:        it.NmID,
			Brand:       it.Brand,
			Status:      it.Status,
		}
	}

	return &models.Order{
		OrderUID:    orderDB.OrderUID,
		TrackNumber: orderDB.TrackNumber,
		Entry:       orderDB.Entry,
		Delivery: models.Delivery{
			Name:    deliveryDB.Name,
			Phone:   deliveryDB.Phone,
			Zip:     deliveryDB.Zip,
			City:    deliveryDB.City,
			Address: deliveryDB.Address,
			Region:  deliveryDB.Region,
			Email:   deliveryDB.Email,
		},
		Payment: models.Payment{
			Transaction:  paymentDB.Transaction,
			RequestID:    paymentDB.RequestID,
			Currency:     paymentDB.Currency,
			Provider:     paymentDB.Provider,
			Amount:       paymentDB.Amount,
			PaymentDt:    paymentDB.PaymentDt,
			Bank:         paymentDB.Bank,
			DeliveryCost: paymentDB.DeliveryCost,
			GoodsTotal:   paymentDB.GoodsTotal,
			CustomFee:    paymentDB.CustomFee,
		},
		Items:             items,
		Locale:            orderDB.Locale,
		InternalSignature: orderDB.InternalSignature,
		CustomerID:        orderDB.CustomerID,
		DeliveryService:   orderDB.DeliveryService,
		Shardkey:          orderDB.Shardkey,
		SmID:              orderDB.SmID,
		DateCreated:       time.Unix(orderDB.DateCreated, 0),
		OofShard:          orderDB.OofShard,
	}
}
