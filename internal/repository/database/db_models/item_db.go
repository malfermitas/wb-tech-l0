package db_models

import (
	"wb-tech-l0/internal/models"

	"gorm.io/gorm"
)

type ItemDB struct {
	gorm.Model
	ChrtID      int
	TrackNumber string
	Price       int
	Rid         string
	Name        string
	Sale        int
	Size        string
	TotalPrice  int
	NmID        int
	Brand       string
	Status      int
	OrderUID    string `gorm:"not null;index"`
}

func ToItemDB(item models.Item, orderUID string) ItemDB {
	return ItemDB{
		ChrtID:      item.ChrtID,
		TrackNumber: item.TrackNumber,
		Price:       item.Price,
		Rid:         item.RID,
		Name:        item.Name,
		Sale:        item.Sale,
		Size:        item.Size,
		TotalPrice:  item.TotalPrice,
		NmID:        item.NmID,
		Brand:       item.Brand,
		Status:      item.Status,
		OrderUID:    orderUID,
	}
}
