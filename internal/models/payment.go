package models

import "gorm.io/gorm"

type Payment struct {
	Transaction  string `json:"transaction" fake:"{uuid}"`
	RequestID    string `json:"request_id" fake:"{regex:[a-zA-Z0-9]{0,10}}"`
	Currency     string `json:"currency" fake:"{randomstring:[USD,RUB,EUR]}"`
	Provider     string `json:"provider" fake:"{randomstring:[wbpay,paypal,stripe]}"`
	Amount       int    `json:"amount" fake:"{number:100,10000}"`
	PaymentDt    int64  `json:"payment_dt" fake:"{number:1609459200,1640995200}"` // 2021-2022 timestamp range
	Bank         string `json:"bank" fake:"{randomstring:[alpha,sber,tinkoff]}"`
	DeliveryCost int    `json:"delivery_cost" fake:"{number:100,2000}"`
	GoodsTotal   int    `json:"goods_total" fake:"{number:50,5000}"`
	CustomFee    int    `json:"custom_fee" fake:"{number:0,100}"`
}

type PaymentDB struct {
	gorm.Model
	Payment
	OrderUID string `gorm:"not null;uniqueIndex" json:"-"`
}
