package models

type Item struct {
	ChrtID      int    `json:"chrt_id" fake:"{number:1000000,9999999}" validate:"required,min=1"`
	TrackNumber string `json:"track_number" fake:"{regex:[A-Z]{10}}" validate:"required,min=5,max=50"`
	Price       int    `json:"price" fake:"{number:100,5000}" validate:"required,min=1"`
	RID         string `json:"rid" fake:"{uuid}" validate:"required,uuid"`
	Name        string `json:"name" fake:"{productname}" validate:"required,min=1,max=200"`
	Sale        int    `json:"sale" fake:"{number:0,50}" validate:"min=0,max=100"`
	Size        string `json:"size" fake:"{number:1,70}" validate:"required,min=1,max=70"`
	TotalPrice  int    `json:"total_price" fake:"{number:100,4500}" validate:"required,min=1"`
	NmID        int    `json:"nm_id" fake:"{number:100000,999999}" validate:"required,min=1"`
	Brand       string `json:"brand" fake:"{company}" validate:"required,min=1,max=100"`
	Status      int    `json:"status" fake:"{number:200,202}" validate:"required,oneof=200 201 202"`
}
