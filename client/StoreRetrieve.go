package client

import (
	"bytes"
	"encoding/json"
	"strings"

	"io/ioutil"

	"net/http"
	"time"

	"github.com/ariefsam/eventsam"
	"github.com/pkg/errors"
)

func init() {
	http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout = time.Second * 45
}

type Eventsam struct {
	server string
}

func NewEventsam(server string) (es Eventsam, err error) {
	server = strings.TrimSuffix(server, "/")
	es = Eventsam{
		server: server,
	}
	return
}

type EventData struct {
	AggregateName string `json:"aggregate_name"`
	AggregateID   string `json:"aggregate_id"`
	EventName     string `json:"event_name"`
	Version       int64  `json:"version"`
	Data          any    `json:"data"`
	TimeMillis    int64  `json:"time_millis"`
}

func (es Eventsam) Store(aggregateID string, aggregateName string, eventName string, version int64, data any) (respEntity eventsam.EventEntity, err error) {

	entity := EventData{
		AggregateID:   aggregateID,
		AggregateName: aggregateName,
		EventName:     eventName,
		Data:          data,
		Version:       version,
		TimeMillis:    time.Now().UnixMilli(),
	}

	payload, err := json.Marshal(entity)
	if err != nil {
		err = errors.Wrap(err, "error marshal to json")
		return
	}
	payloadReader := bytes.NewReader(payload)
	res, err := http.Post(es.server+"/store", "application/json", payloadReader)
	if err != nil {
		err = errors.Wrap(err, "error post to server eventsam")
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = errors.Wrap(err, "error read response body from eventsam")
		return
	}

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
	filter := map[string]any{
		"aggregate_id":   aggregateID,
		"aggregate_name": aggregateName,
		"after_version":  sinceVersion,
	}
	jsonFilter, err := json.Marshal(filter)
	if err != nil {
		return
	}
	filterReader := bytes.NewReader(jsonFilter)

	res, err := http.Post(es.server+"/retrieve", "application/json", filterReader)
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

func (es *Eventsam) FetchAggregateEvent(aggregateName string, afterID, limit int) (events []eventsam.EventEntity, err error) {
	filter := map[string]any{
		"aggregate_name": aggregateName,
		"after_id":       afterID,
		"limit":          limit,
	}

	jsonFilter, err := json.Marshal(filter)
	if err != nil {
		return
	}
	filterReader := bytes.NewReader(jsonFilter)

	res, err := http.Post(es.server+"/fetch-aggregate-event", "application/json", filterReader)
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
