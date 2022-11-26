package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/PT-Jojonomic-Indonesia/microkit/response"
	"github.com/ariefsam/eventsam"
	"github.com/ariefsam/eventsam/client"
)

var cond *sync.Cond
var condLock sync.Mutex
var condAggregate map[string]*sync.Cond

func init() {
	condAggregate = make(map[string]*sync.Cond)
}

type EventData struct {
	AggregateName string `json:"aggregate_name"`
	AggregateID   string `json:"aggregate_id"`
	EventName     string `json:"event_name"`
	Version       int64  `json:"version"`
	Data          any    `json:"data"`
	TimeMillis    int64  `json:"time_millis"`
}

func StoreHandler(w http.ResponseWriter, r *http.Request) {

	data := EventData{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		response.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	var event eventsam.EventEntity
	if isSlave {
		esamClient, err := client.NewEventsam(os.Getenv("MASTER_ADDRESS"))
		if err != nil {
			response.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}
		event, err = esamClient.Store(data.AggregateID, data.AggregateName, data.EventName, data.Version, data.Data)
		if err != nil {
			log.Println(err)
			response.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}
		goto response
	}

	event, err = esam.Store(data.AggregateID, data.AggregateName, data.EventName, data.Version, data.Data)
	if err != nil {
		response.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	go func() {
		defer func() {
			recover()
		}()
		cond.L.Lock()
		cond.Broadcast()
		cond.L.Unlock()
	}()
	go func() {
		defer func() {
			recover()
		}()
		condA := getCondAggregate(data.AggregateName)
		condA.L.Lock()
		condA.Broadcast()
		condA.L.Unlock()
	}()
response:
	dataResp := map[string]any{
		"message": "success",
		"data":    event,
	}
	response.JSON(w, dataResp, http.StatusOK)
}
