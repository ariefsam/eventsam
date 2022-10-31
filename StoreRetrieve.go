package eventsam

import (
	"encoding/json"
	"eventsam/idgenerator"
	"time"
)

func (es Eventsam) Store(aggregateID string, aggregateName string, eventName string, version int64, data any) (err error) {
	tmp, _ := json.Marshal(data)
	dataString := string(tmp)
	entity := EventEntity{
		AggregateID:   aggregateID,
		AggregateName: aggregateName,
		EventID:       idgenerator.Generate(),
		EventName:     eventName,
		Data:          dataString,
		Version:       version,
		TimeMillis:    time.Now().UnixMilli(),
	}
	err = es.db.Save(&entity).Error
	return
}

func (es *Eventsam) Retrieve(aggregateID string, aggregateName string, version int) (events []EventEntity, err error) {
	err = es.db.Where("aggregate_id = ? AND aggregate_name = ? AND version >= ? ", aggregateID, aggregateName, version).Find(&events).Error
	return
}

func (es *EventEntity) DataToStruct(data any) (err error) {
	err = json.Unmarshal([]byte(es.Data), data)
	return
}
