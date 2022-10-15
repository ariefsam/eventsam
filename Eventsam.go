package eventsam

import (
	"gorm.io/gorm"
)

type EventEntity struct {
	ID            uint   `gorm:"primarykey"`
	EventID       string `gorm:"uniqueIndex"`
	AggregateName string `gorm:"index:,unique,composite:idx_aggregate_version;index"`
	AggregateID   string `gorm:"index:,unique,composite:idx_aggregate_version;index"`
	EventName     string `gorm:"index"`
	Version       int64  `gorm:"index:,unique,composite:idx_aggregate_version;index"`
	Data          string
	TimeMillis    int64
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
