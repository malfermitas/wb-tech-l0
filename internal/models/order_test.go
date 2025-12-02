package models_test

import (
	"encoding/json"
	"testing"
	"time"
	"wb-tech-l0/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderJSONSerialization(t *testing.T) {
	order := models.Order{
		OrderUID:          "test-uid",
		TrackNumber:       "TRACK123",
		Entry:             "WBIL",
		Locale:            "en",
		InternalSignature: "test-signature",
		CustomerID:        "customer-123",
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
			Transaction:  "test-transaction",
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
				TrackNumber: "TRACK123",
				Price:       500,
				RID:         "test-rid",
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

	// Test JSON marshaling
	jsonData, err := json.Marshal(order)
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Test JSON unmarshaling
	var unmarshaled models.Order
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, order.OrderUID, unmarshaled.OrderUID)
	assert.Equal(t, order.TrackNumber, unmarshaled.TrackNumber)
	assert.Equal(t, order.CustomerID, unmarshaled.CustomerID)
	assert.Equal(t, len(order.Items), len(unmarshaled.Items))
}
