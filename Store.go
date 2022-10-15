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
