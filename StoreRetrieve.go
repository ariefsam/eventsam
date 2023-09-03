package eventsam

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"

	"time"

	"github.com/ariefsam/eventsam/idgenerator"
)

func (es Eventsam) Store(aggregateID string, aggregateName string, eventName string, version int64, data any) (entity EventEntity, err error) {
	if aggregateID == "" {
		err = errors.New("aggregate_id is empty")
		return
	}
	if aggregateName == "" {
		err = errors.New("aggregate_name is empty")
		return
	}
	if eventName == "" {
		err = errors.New("event_name is empty")
		return
	}
	if version < 1 {
		err = errors.New("minimal version is 1")
		return
	}
	if data == nil {
		err = errors.New("data is empty")
		return
	}

	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	err = jsonEncoder.Encode(data)
	if err != nil {
		log.Println(err)
		return
	}

	dataString := bf.String()

	if version > 1 {

		oldEvent := []EventEntity{}
		err = es.db.Limit(1).Order("id DESC").Where("aggregate_id = ? AND aggregate_name = ? ", aggregateID, aggregateName).Find(&oldEvent).Error
		if err != nil {
			return
		}

		if len(oldEvent) == 0 {
			if version > 1 {
				err = errors.New("please start with version 1")
				return
			}
		}

		if len(oldEvent) > 0 {
			diffVersion := version - oldEvent[0].Version
			if diffVersion == 0 {
				err = errors.New("duplicate version")
				return
			}
			if version-oldEvent[0].Version > 1 {
				err = errors.New("version is not sequential")
				return
			}
		}
	}

	entity = EventEntity{
		AggregateID:   aggregateID,
		AggregateName: aggregateName,
		EventID:       idgenerator.Generate(),
		EventName:     eventName,
		Data:          dataString,
		Version:       version,
		TimeMillis:    time.Now().UnixMilli(),
	}
	err = es.db.Save(&entity).Error
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: event_entities.aggregate_name, event_entities.aggregate_id, event_entities.version" {
			err = errors.New("duplicate version")
			return
		}
	}
	return
}

func (es *Eventsam) Retrieve(aggregateID string, aggregateName string, afterVersion int) (events []EventEntity, err error) {
	err = es.db.Where("aggregate_id = ? AND aggregate_name = ? AND version > ? ", aggregateID, aggregateName, afterVersion).Order("id asc").Find(&events).Error
	return
}

func (es *Eventsam) FetchAllEvent(afterID, limit int) (events []EventEntity, err error) {
	err = es.db.Where("id > ? ", afterID).Order("id asc").Limit(limit).Find(&events).Error
	return
}

func (es *Eventsam) FetchAggregateEvent(aggregateName string, afterID, limit int) (events []EventEntity, err error) {
	err = es.db.Where("id > ? AND aggregate_name = ? ", afterID, aggregateName).Order("id asc").Limit(limit).Find(&events).Error
	return
}
