package eventsam

import (
	"gorm.io/gorm"
)

type EventEntity struct {
	ID            uint   `gorm:"primarykey"`
	EventID       string `gorm:"uniqueIndex"`
	AggregateType string `gorm:"index:,unique,composite:idx_aggregate_version;index"`
	AggregateID   string `gorm:"index:,unique,composite:idx_aggregate_version;index"`
	EventName     string `gorm:"index:,unique,composite:idx_aggregate_version;index"`
	Version       int    `gorm:"index:,unique,composite:idx_aggregate_version;index"`
	Data          string
	Timestamp     int
}

type Eventsam struct {
	db *gorm.DB
}

func NewEventsam(db *gorm.DB) (es Eventsam, err error) {
	es = Eventsam{
		db: db,
	}
	es.db.AutoMigrate(&EventEntity{})
	return
}
