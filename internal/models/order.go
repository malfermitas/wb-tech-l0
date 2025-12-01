package models

import (
	"time"
)

type Order struct {
	OrderUID          string    `json:"order_uid" fake:"{uuid}"`
	TrackNumber       string    `json:"track_number" fake:"{regex:[A-Z]{10}}"`
	Entry             string    `json:"entry" fake:"{regex:[A-Z]{4}}"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items" fakesize:"1,5"`
	Locale            string    `json:"locale" fake:"{randomstring:[en,ru,es]}"`
	InternalSignature string    `json:"internal_signature" fake:"{regex:[a-zA-Z0-9]{0,20}}"`
	CustomerID        string    `json:"customer_id" fake:"{uuid}"`
	DeliveryService   string    `json:"delivery_service" fake:"{randomstring:[meest,ups,fedex,dhl]}"`
	Shardkey          string    `json:"shardkey" fake:"{number:1,9}"`
	SmID              int       `json:"sm_id" fake:"{number:1,100}"`
	DateCreated       time.Time `json:"date_created" fake:"{date}"`
	OofShard          string    `json:"oof_shard" fake:"{number:1,9}"`
}
