package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/PT-Jojonomic-Indonesia/microkit/response"
)

type FetchAllEventInput struct {
	WaitTimeMillisIfEmpty int64 `json:"wait_time_millis_if_empty"`
	AfterID               int   `json:"after_id"`
	Limit                 int   `json:"limit"`
}

func FetchAllEventHandler(w http.ResponseWriter, r *http.Request) {
	dataInput := FetchAllEventInput{}
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
			defer func() {
				if r := recover(); r != nil {
					log.Println("Recovered in f", r)
				}
			}()
			cond.L.Lock()
			cond.Wait()
			cond.L.Unlock()
			c <- true
		}()
		select {
		case <-c:
		case <-time.After(time.Duration(dataInput.WaitTimeMillisIfEmpty) * time.Millisecond):
		}

		close(c)
		events, err = esam.FetchAllEvent(dataInput.AfterID, dataInput.Limit)
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
