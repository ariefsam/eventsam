package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"

	"log"
	"net/http"
	"time"

	"github.com/ariefsam/eventsam"
	"github.com/ariefsam/eventsam/idgenerator"
)

func init() {
	http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout = time.Second * 45
}

type Eventsam struct {
	server string
}

func NewEventsam(server string) (es Eventsam, err error) {
	es = Eventsam{
		server: server,
	}
	return
}

func (es Eventsam) Store(aggregateID string, aggregateName string, eventName string, version int64, data any) (respEntity eventsam.EventEntity, err error) {
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
	payload, _ := json.Marshal(entity)
	payloadReader := bytes.NewReader(payload)
	res, err := http.Post(es.server+"/store", "application/json", payloadReader)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	log.Println("body ", string(body))
	resp := struct {
		Error   bool                 `json:"error"`
		Message string               `json:"message"`
		Data    eventsam.EventEntity `json:"data"`
	}{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return
	}
	if resp.Error {
		err = errors.New(resp.Message)
		return
	}
	respEntity = resp.Data
	return
}

func (es *Eventsam) Retrieve(aggregateID string, aggregateName string, sinceVersion int) (events []eventsam.EventEntity, err error) {

	return
}

func (es *Eventsam) FetchAllEvent(afterID, limit int) (events []eventsam.EventEntity, err error) {
	filter := map[string]any{
		"after_id": afterID,
		"limit":    limit,
	}
	jsonFilter, err := json.Marshal(filter)
	if err != nil {
		return
	}
	filterReader := bytes.NewReader(jsonFilter)

	res, err := http.Post(es.server+"/fetch-all-event", "application/json", filterReader)
	if err != nil {
		return
	}
	defer res.Body.Close()
	eventResponse := struct {
		Data    []eventsam.EventEntity `json:"data"`
		Error   bool                   `json:"error"`
		Message string                 `json:"message"`
	}{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &eventResponse)
	if err != nil {
		return
	}
	events = eventResponse.Data
	return
}
