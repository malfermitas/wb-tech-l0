package models

type Payment struct {
	Transaction  string `json:"transaction" fake:"{uuid}" validate:"required,uuid"`
	RequestID    string `json:"request_id" fake:"{uuid}" validate:"omitempty,uuid"`
	Currency     string `json:"currency" fake:"{currencyshort}" validate:"required,len=3"`
	Provider     string `json:"provider" fake:"{randomstring:[wbpay,paypal,stripe]}" validate:"required,oneof=wbpay paypal stripe"`
	Amount       int    `json:"amount" fake:"{number:1,10000}" validate:"required,min=1"`
	PaymentDt    int64  `json:"payment_dt" fake:"{unixtime}" validate:"required"`
	Bank         string `json:"bank" fake:"{randomstring:[alpha,sberbank,tinkoff]}" validate:"required,min=2,max=50"`
	DeliveryCost int    `json:"delivery_cost" fake:"{number:100,1000}" validate:"min=0"`
	GoodsTotal   int    `json:"goods_total" fake:"{number:1000,9000}" validate:"required,min=1"`
	CustomFee    int    `json:"custom_fee" fake:"{number:0,500}" validate:"min=0"`
}
