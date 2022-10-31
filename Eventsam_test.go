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

	//error version
	err = esam.Store(aggregateID, "item", "item_received", 0, received)
	assert.Error(t, err)

	err = esam.Store(aggregateID, "item", "item_received", 1, received)
	assert.NoError(t, err)

	events, err := esam.Retrieve(aggregateID, "item", 0)
	assert.NoError(t, err)

	item := Item{}
	item.HandleEvents(events)

}

type Item struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

func (i *Item) HandleEvents(event []eventsam.EventEntity) (err error) {

	return
}

func (i *Item) HandleEvent(event eventsam.EventEntity) (err error) {
	switch event.EventName {
	case "item_purchased":
		err = i.handleItemPurchased(event)
	case "item_received":
		err = i.handleItemReceived(event)
	}
	return
}

func (i *Item) handleItemPurchased(event eventsam.EventEntity) (err error) {
	purchase := struct {
		SKU          string  `json:"sku"`
		PricePerItem float64 `json:"price_per_item"`
		Quantity     int     `json:"quantity"`
		PurchaseID   string  `json:"purchase_id"`
	}{}
	err = event.DataToStruct(&purchase)
	if err != nil {
		return
	}
	i.ID = event.AggregateID
	i.Quantity = purchase.Quantity
	return
}

func (i *Item) handleItemReceived(event eventsam.EventEntity) (err error) {
	received := struct {
		PurchaseID string `json:"purchase_id"`
		Quantity   int    `json:"quantity"`
	}{}
	err = event.DataToStruct(&received)
	if err != nil {
		return
	}
	i.Quantity += received.Quantity
	return
}
