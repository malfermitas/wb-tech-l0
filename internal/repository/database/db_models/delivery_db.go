package db_models

import (
	"wb-tech-l0/internal/models"

	"gorm.io/gorm"
)

type DeliveryDB struct {
	gorm.Model
	Name    string
	Phone   string
	Zip     string
	City    string
	Address string
	Region  string
	Email   string
}

func ToDeliveryDB(d models.Delivery) DeliveryDB {
	return DeliveryDB{
		Name:    d.Name,
		Phone:   d.Phone,
		Zip:     d.Zip,
		City:    d.City,
		Address: d.Address,
		Region:  d.Region,
		Email:   d.Email,
	}
}
