package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"log"
	"net/http"
	"time"

	"github.com/ariefsam/eventsam"
	"github.com/ariefsam/eventsam/idgenerator"
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
	filter := map[string]any{
		"after_id": afterID,
		"limit":    limit,
	}
	jsonFilter, err := json.Marshal(filter)
	if err != nil {
		return
	}
	filterReader := bytes.NewReader(jsonFilter)
	http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout = time.Second * 45
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
	log.Println(string(body))
	err = json.Unmarshal(body, &eventResponse)
	if err != nil {
		return
	}
	events = eventResponse.Data
	return
}
