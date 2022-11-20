package client

import (
	"encoding/json"
	"eventsam"
	"eventsam/idgenerator"
	"log"
	"time"
)

type Eventsam struct {
	server string
}

func NewEventsam(server string) (es Eventsam, err error) {
	es = Eventsam{
		server: server,
	}
	return
}

func (es Eventsam) Store(aggregateID string, aggregateName string, eventName string, version int64, data any) (err error) {
	tmp, _ := json.Marshal(data)
	dataString := string(tmp)
	entity := eventsam.EventEntity{
		AggregateID:   aggregateID,
		AggregateName: aggregateName,
		EventID:       idgenerator.Generate(),
		EventName:     eventName,
		Data:          dataString,
		Version:       version,
		TimeMillis:    time.Now().UnixMilli(),
	}
	// err = es.db.Save(&entity).Error
	log.Println(entity)
	return
}

func (es *Eventsam) Retrieve(aggregateID string, aggregateName string, sinceVersion int) (events []eventsam.EventEntity, err error) {

	return
}

func (es *Eventsam) FetchAllEvent(afterID, limit int) (events []eventsam.EventEntity, err error) {

	return
}
