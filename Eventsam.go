package eventsam

import (
	"encoding/json"

	"gorm.io/gorm"
)

type EventEntity struct {
	ID            uint   `gorm:"primarykey" json:"id"`
	EventID       string `gorm:"uniqueIndex" json:"event_id"`
	AggregateName string `gorm:"index:,unique,composite:idx_aggregate_version;index" json:"aggregate_name"`
	AggregateID   string `gorm:"index:,unique,composite:idx_aggregate_version;index" json:"aggregate_id"`
	EventName     string `gorm:"index" json:"event_name"`
	Version       int64  `gorm:"index:,unique,composite:idx_aggregate_version;index" json:"version"`
	Data          string `json:"data"`
	TimeMillis    int64  `json:"time_millis"`
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

func (es *EventEntity) DataToStruct(data any) (err error) {
	err = json.Unmarshal([]byte(es.Data), data)
	return
}
