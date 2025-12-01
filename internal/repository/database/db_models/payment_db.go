package db_models

import (
	"wb-tech-l0/internal/models"

	"gorm.io/gorm"
)

type PaymentDB struct {
	gorm.Model
	OrderUID     string `gorm:"index;not null"`
	Transaction  string
	RequestID    string
	Currency     string
	Provider     string
	Amount       int
	PaymentDt    int64
	Bank         string
	DeliveryCost int
	GoodsTotal   int
	CustomFee    int
}

func ToPaymentDB(o *models.Order) PaymentDB {
	p := o.Payment
	return PaymentDB{
		OrderUID:     o.OrderUID,
		Transaction:  p.Transaction,
		RequestID:    p.RequestID,
		Currency:     p.Currency,
		Provider:     p.Provider,
		Amount:       p.Amount,
		PaymentDt:    p.PaymentDt,
		Bank:         p.Bank,
		DeliveryCost: p.DeliveryCost,
		GoodsTotal:   p.GoodsTotal,
		CustomFee:    p.CustomFee,
	}
}
