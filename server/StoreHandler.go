package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/PT-Jojonomic-Indonesia/microkit/response"
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

	event, err := esam.Store(data.AggregateID, data.AggregateName, data.EventName, data.Version, data.Data)
	if err != nil {
		log.Println(err)
		response.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	dataResp := map[string]any{
		"message": "success",
		"data":    event,
	}
	response.JSON(w, dataResp, http.StatusOK)
}
