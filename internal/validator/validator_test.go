package validator_test

import (
	"testing"
	"time"
	"wb-tech-l0/internal/models"
	"wb-tech-l0/internal/validator"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidator_ValidateOrder_Valid(t *testing.T) {
	v := validator.NewValidator()

	order := models.Order{
		OrderUID:          uuid.New().String(),
		TrackNumber:       "TRACK12345",
		Entry:             "WBIL",
		Locale:            "en",
		InternalSignature: "signature",
		CustomerID:        uuid.New().String(),
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       time.Now(),
		OofShard:          "1",
		Delivery: models.Delivery{
			Name:    "John Doe",
			Phone:   "+1234567890",
			Zip:     "12345",
			City:    "Test City",
			Address: "123 Test St",
			Region:  "Test Region",
			Email:   "test@example.com",
		},
		Payment: models.Payment{
			Transaction:  uuid.New().String(),
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1000,
			PaymentDt:    time.Now().Unix(),
			Bank:         "test-bank",
			DeliveryCost: 100,
			GoodsTotal:   900,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtID:      1234567,
				TrackNumber: "TRACK12345",
				Price:       500,
				RID:         uuid.New().String(),
				Name:        "Test Item",
				Sale:        10,
				Size:        "M",
				TotalPrice:  450,
				NmID:        123456,
				Brand:       "Test Brand",
				Status:      200,
			},
		},
	}

	err := v.Validate(&order)
	assert.NoError(t, err)
}

func TestValidator_ValidateOrder_Invalid(t *testing.T) {
	v := validator.NewValidator()

	tests := []struct {
		name  string
		order models.Order
	}{
		{
			name: "missing order_uid",
			order: models.Order{
				TrackNumber: "TRACK12345",
			},
		},
		{
			name: "invalid locale",
			order: models.Order{
				OrderUID:        uuid.New().String(),
				TrackNumber:     "TRACK12345",
				Entry:           "WBIL",
				Locale:          "invalid",
				CustomerID:      uuid.New().String(),
				DeliveryService: "meest",
				Shardkey:        "9",
				SmID:            99,
				DateCreated:     time.Now(),
				OofShard:        "1",
			},
		},
		{
			name: "invalid email in delivery",
			order: models.Order{
				OrderUID:        uuid.New().String(),
				TrackNumber:     "TRACK12345",
				Entry:           "WBIL",
				Locale:          "en",
				CustomerID:      uuid.New().String(),
				DeliveryService: "meest",
				Shardkey:        "9",
				SmID:            99,
				DateCreated:     time.Now(),
				OofShard:        "1",
				Delivery: models.Delivery{
					Name:    "John Doe",
					Phone:   "+1234567890",
					Zip:     "12345",
					City:    "Test City",
					Address: "123 Test St",
					Region:  "Test Region",
					Email:   "invalid-email",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(&tt.order)
			assert.Error(t, err)
		})
	}
}
