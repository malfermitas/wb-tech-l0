package validator_test

import (
	"testing"
	"time"
	"wb-tech-l0/internal/models"
	vpkg "wb-tech-l0/internal/validator"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// --- Helpers for valid fixtures ---
func validDelivery() models.Delivery {
	return models.Delivery{
		Name:    "John Doe",
		Phone:   "+12345678901",
		Zip:     "12345",
		City:    "City",
		Address: "123 Test Street",
		Region:  "Region",
		Email:   "john@example.com",
	}
}

func validPayment() models.Payment {
	return models.Payment{
		Transaction:  uuid.New().String(),
		RequestID:    "",
		Currency:     "USD",
		Provider:     "wbpay",
		Amount:       100,
		PaymentDt:    time.Now().Unix(),
		Bank:         "alpha",
		DeliveryCost: 10,
		GoodsTotal:   90,
		CustomFee:    0,
	}
}

func validItem(track string) models.Item {
	return models.Item{
		ChrtID:      1234567,
		TrackNumber: track,
		Price:       500,
		RID:         uuid.New().String(),
		Name:        "Item",
		Sale:        0,
		Size:        "M",
		TotalPrice:  500,
		NmID:        123456,
		Brand:       "Brand",
		Status:      200,
	}
}

func validOrder() models.Order {
	track := "TRACKABCDE"
	return models.Order{
		OrderUID:          uuid.New().String(),
		TrackNumber:       track,
		Entry:             "WBIL",
		Delivery:          validDelivery(),
		Payment:           validPayment(),
		Items:             []models.Item{validItem(track)},
		Locale:            "en",
		InternalSignature: "sig",
		CustomerID:        uuid.New().String(),
		DeliveryService:   "meest",
		Shardkey:          "1",
		SmID:              1,
		DateCreated:       time.Now(),
		OofShard:          "1",
	}
}

// --- Delivery ---
func TestDeliveryValidation(t *testing.T) {
	v := vpkg.NewValidator()

	t.Run("valid", func(t *testing.T) {
		d := validDelivery()
		assert.NoError(t, v.Validate(d))
	})

	t.Run("invalid email", func(t *testing.T) {
		d := validDelivery()
		d.Email = "not-an-email"
		assert.Error(t, v.Validate(d))
	})

	t.Run("invalid phone not e164", func(t *testing.T) {
		d := validDelivery()
		d.Phone = "123456" // missing leading + and country code
		assert.Error(t, v.Validate(d))
	})
}

// --- Payment ---
func TestPaymentValidation(t *testing.T) {
	v := vpkg.NewValidator()

	t.Run("valid", func(t *testing.T) {
		p := validPayment()
		assert.NoError(t, v.Validate(p))
	})

	t.Run("invalid provider enum", func(t *testing.T) {
		p := validPayment()
		p.Provider = "unknown"
		assert.Error(t, v.Validate(p))
	})

	t.Run("invalid currency length", func(t *testing.T) {
		p := validPayment()
		p.Currency = "US"
		assert.Error(t, v.Validate(p))
	})

	t.Run("amount must be >=1", func(t *testing.T) {
		p := validPayment()
		p.Amount = 0
		assert.Error(t, v.Validate(p))
	})

	t.Run("negative custom fee not allowed", func(t *testing.T) {
		p := validPayment()
		p.CustomFee = -1
		assert.Error(t, v.Validate(p))
	})
}

// --- Item ---
func TestItemValidation(t *testing.T) {
	v := vpkg.NewValidator()
	track := "TRACKABCDE"

	t.Run("valid", func(t *testing.T) {
		it := validItem(track)
		assert.NoError(t, v.Validate(it))
	})

	t.Run("rid must be uuid", func(t *testing.T) {
		it := validItem(track)
		it.RID = "not-uuid"
		assert.Error(t, v.Validate(it))
	})

	t.Run("status must be in enum", func(t *testing.T) {
		it := validItem(track)
		it.Status = 203
		assert.Error(t, v.Validate(it))
	})

	t.Run("total_price must be >=1", func(t *testing.T) {
		it := validItem(track)
		it.TotalPrice = 0
		assert.Error(t, v.Validate(it))
	})
}

// --- Order nested ---
func TestOrderValidation_ItemsAndNested(t *testing.T) {
	v := vpkg.NewValidator()

	t.Run("valid order", func(t *testing.T) {
		o := validOrder()
		assert.NoError(t, v.Validate(o))
	})

	t.Run("empty items invalid", func(t *testing.T) {
		o := validOrder()
		o.Items = nil
		assert.Error(t, v.Validate(o))
		o.Items = []models.Item{}
		assert.Error(t, v.Validate(o))
	})

	t.Run("invalid nested item", func(t *testing.T) {
		o := validOrder()
		o.Items[0].RID = "bad"
		assert.Error(t, v.Validate(o))
	})

	t.Run("invalid uuids in order fields", func(t *testing.T) {
		o := validOrder()
		o.OrderUID = "bad"
		o.CustomerID = "bad"
		assert.Error(t, v.Validate(o))
	})

	t.Run("sm_id must be >=1", func(t *testing.T) {
		o := validOrder()
		o.SmID = 0
		assert.Error(t, v.Validate(o))
	})

	t.Run("track number too short", func(t *testing.T) {
		o := validOrder()
		o.TrackNumber = "A1"
		assert.Error(t, v.Validate(o))
	})
}
