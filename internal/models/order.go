package models

import (
	"time"
)

type Order struct {
	OrderUID          string    `json:"order_uid" fake:"{uuid}" validate:"required,uuid"`
	TrackNumber       string    `json:"track_number" fake:"{regex:[A-Z]{10}}" validate:"required,min=5,max=50"`
	Entry             string    `json:"entry" fake:"{regex:[A-Z]{4}}" validate:"required,min=2,max=10"`
	Delivery          Delivery  `json:"delivery" validate:"required"`
	Payment           Payment   `json:"payment" validate:"required"`
	Items             []Item    `json:"items" fakesize:"1,5" validate:"required,min=1,dive"`
	Locale            string    `json:"locale" fake:"{randomstring:[en,ru,es]}" validate:"required,oneof=en ru es"`
	InternalSignature string    `json:"internal_signature" fake:"{regex:[a-zA-Z0-9]{0,20}}" validate:"max=50"`
	CustomerID        string    `json:"customer_id" fake:"{uuid}" validate:"required,uuid"`
	DeliveryService   string    `json:"delivery_service" fake:"{randomstring:[meest,ups,fedex,dhl]}" validate:"required,oneof=meest ups fedex dhl"`
	Shardkey          string    `json:"shardkey" fake:"{number:1,9}" validate:"required,min=1"`
	SmID              int       `json:"sm_id" fake:"{number:1,100}" validate:"required,min=1"`
	DateCreated       time.Time `json:"date_created" fake:"{date}" validate:"required"`
	OofShard          string    `json:"oof_shard" fake:"{number:1,9}" validate:"required,min=1"`
}
