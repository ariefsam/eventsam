package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/PT-Jojonomic-Indonesia/microkit/response"
	"github.com/ariefsam/eventsam"
	"github.com/ariefsam/eventsam/client"
)

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
		cond.L.Lock()
		cond.Broadcast()
		cond.L.Unlock()
	}()
response:
	dataResp := map[string]any{
		"message": "success",
		"data":    event,
	}
	response.JSON(w, dataResp, http.StatusOK)
}
