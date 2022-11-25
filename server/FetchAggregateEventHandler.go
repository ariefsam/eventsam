package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/PT-Jojonomic-Indonesia/microkit/response"
)

type FetchAggregateEventInput struct {
	AggregateName string `json:"aggregate_name"`
	AfterID       int    `json:"after_id"`
	Limit         int    `json:"limit"`
}

func FetchAggregateEventHandler(w http.ResponseWriter, r *http.Request) {
	dataInput := FetchAggregateEventInput{}
	err := json.NewDecoder(r.Body).Decode(&dataInput)
	if err != nil {
		response.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	events, err := esam.FetchAllEvent(dataInput.AfterID, dataInput.Limit)
	if err != nil {
		response.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	if len(events) == 0 {
		c := make(chan bool)
		go func() {
			defer recover()
			cond.L.Lock()
			cond.Wait()
			cond.L.Unlock()
			c <- true
		}()
		select {
		case <-c:
		case <-time.After(25 * time.Second):
		}

		close(c)
		events, err = esam.FetchAggregateEvent(dataInput.AggregateName, dataInput.AfterID, dataInput.Limit)
		if err != nil {
			response.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}
	}
	dataResp := map[string]any{
		"message": "success",
		"data":    events,
		"error":   false,
	}
	response.JSON(w, dataResp, http.StatusOK)
}
