package database_test

import (
	"fmt"
	"testing"
	"time"

	"wb-tech-l0/internal/models"
	"wb-tech-l0/internal/repository/cache"
	dbpkg "wb-tech-l0/internal/repository/database"
	"wb-tech-l0/internal/repository/database/db_models"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestOrder(uid string) *models.Order {
	if uid == "" {
		uid = "uid-1"
	}
	return &models.Order{
		OrderUID:        uid,
		TrackNumber:     "ABCDEFGHJK",
		Entry:           "WBIL",
		Locale:          "en",
		CustomerID:      "11111111-1111-1111-1111-111111111111",
		DeliveryService: "meest",
		Shardkey:        "1",
		SmID:            1,
		DateCreated:     time.Now(),
		OofShard:        "1",
		Delivery:        models.Delivery{Name: "John", Phone: "+123", Zip: "12345", City: "City", Address: "Street 1", Region: "Reg", Email: "john@example.com"},
		Payment:         models.Payment{Transaction: "tx1", Currency: "USD", Provider: "wbpay", Amount: 100, PaymentDt: time.Now().Unix(), Bank: "bank", DeliveryCost: 10, GoodsTotal: 90, CustomFee: 0},
		Items:           []models.Item{{ChrtID: 1, TrackNumber: "ABCDEFGHJK", Price: 100, RID: "rid", Name: "Item", Sale: 0, Size: "M", TotalPrice: 100, NmID: 1, Brand: "brand", Status: 200}},
	}
}

func newTestDB(t *testing.T) (*dbpkg.DB, func()) {
	t.Helper()

	// SQLite in-memory for GORM
	// Use a unique DSN per test to avoid cross-test interference when running the whole suite.
	dsn := fmt.Sprintf("file:orderrepo_%d?mode=memory&cache=shared", time.Now().UnixNano())
	gdb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	// Run migrations for required tables
	err = gdb.AutoMigrate(&db_models.DeliveryDB{}, &db_models.PaymentDB{}, &db_models.OrderDB{}, &db_models.ItemDB{})
	require.NoError(t, err)

	// Miniredis for cache
	mr, err := miniredis.Run()
	require.NoError(t, err)

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	oc := cache.NewOrderCache(rdb, time.Hour)

	db := &dbpkg.DB{Conn: gdb, Cache: oc}
	cleanup := func() {
		rdb.Close()
		mr.Close()
	}
	return db, cleanup
}

func TestOrderRepository_SaveAndGetOrder(t *testing.T) {
	db, cleanup := newTestDB(t)
	defer cleanup()

	order := newTestOrder("uid-save-1")
	err := db.SaveOrder(order)
	require.NoError(t, err)

	// Fetch back
	got, err := db.GetOrder(order.OrderUID)
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, order.OrderUID, got.OrderUID)
	assert.Equal(t, order.Delivery.Email, got.Delivery.Email)
	require.Len(t, got.Items, 1)
	assert.Equal(t, order.Items[0].RID, got.Items[0].RID)
}

func TestOrderRepository_GetOrderCount(t *testing.T) {
	db, cleanup := newTestDB(t)
	defer cleanup()

	// no records yet
	cnt, err := db.GetOrderCount()
	require.NoError(t, err)
	assert.Equal(t, int64(0), cnt)

	// add one
	require.NoError(t, db.SaveOrder(newTestOrder("uid-count-1")))
	cnt, err = db.GetOrderCount()
	require.NoError(t, err)
	assert.Equal(t, int64(1), cnt)
}

func TestOrderRepository_LoadAllOrdersToCache(t *testing.T) {
	db, cleanup := newTestDB(t)
	defer cleanup()

	// populate two orders
	require.NoError(t, db.SaveOrder(newTestOrder("uid-cache-1")))
	require.NoError(t, db.SaveOrder(newTestOrder("uid-cache-2")))

	// Clear redis by recreating client through the same addr is complex; rely on method behavior filling cache
	err := db.LoadAllOrdersToCache()
	require.NoError(t, err)

	// We expect cache to have at least 2 keys
	size := db.CacheSize()
	assert.GreaterOrEqual(t, size, 2)
}
