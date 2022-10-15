package eventsam_test

import (
	"eventsam"
	"eventsam/idgenerator"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestNewEventsam(t *testing.T) {

	filepath := "test_db.db"
	logService := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,         // Disable color
		})
	db, err := gorm.Open(sqlite.Open(filepath), &gorm.Config{
		Logger: logService,
	})
	assert.NoError(t, err)
	// See "Important settings" section.

	esam, err := eventsam.NewEventsam(db)
	if err != nil {
		t.Error(err)

		return
	}
	purchased := struct {
		SKU          string  `json:"sku"`
		PricePerItem float64 `json:"price_per_item"`
		Quantity     int     `json:"quantity"`
		PurchaseID   string  `json:"purchase_id"`
	}{
		SKU:          "sku001",
		PricePerItem: 500.0,
		Quantity:     20,
		PurchaseID:   idgenerator.Generate(),
	}

	aggregateID := idgenerator.Generate()
	err = esam.Store(aggregateID, "item", "item_purchased", 0, purchased)
	assert.NoError(t, err)

	err = esam.Store(aggregateID, "item", "item_purchased", 0, purchased)
	assert.Error(t, err)

	received := struct {
		PurchaseID string `json:"purchase_id"`
		Quantity   int    `json:"quantity"`
	}{
		PurchaseID: purchased.PurchaseID,
		Quantity:   20,
	}
	err = esam.Store(aggregateID, "item", "item_received", 0, received)
	assert.Error(t, err)

	err = esam.Store(aggregateID, "item", "item_received", 1, received)
	assert.NoError(t, err)
}
