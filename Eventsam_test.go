package eventsam_test

import (
	"eventsam"
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
		SKU string `json:"sku"`
	}{
		SKU: "sku001",
	}
	err = esam.Store("a001", "item", "item_purchased", purchased)
	assert.NoError(t, err)
}
