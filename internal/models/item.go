package models

import "gorm.io/gorm"

type Item struct {
	ChrtID      int    `json:"chrt_id" fake:"{number:1000000,9999999}"`
	TrackNumber string `json:"track_number" fake:"{regex:[A-Z]{10}}"`
	Price       int    `json:"price" fake:"{number:100,5000}"`
	Rid         string `json:"rid" fake:"{uuid}"`
	Name        string `json:"name" fake:"{productname}"`
	Sale        int    `json:"sale" fake:"{number:0,50}"`
	Size        string `json:"size" fake:"{number:0,5}"`
	TotalPrice  int    `json:"total_price" fake:"{number:100,5000}"`
	NmID        int    `json:"nm_id" fake:"{number:1000000,9999999}"`
	Brand       string `json:"brand" fake:"{company}"`
	Status      int    `json:"status" fake:"{number:100,400}"`
}

type ItemDB struct {
	gorm.Model
	Item
	OrderUID string `gorm:"not null;index" json:"-"`
}
