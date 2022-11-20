package main

import (
	"encoding/json"
	"net/http"

	"github.com/PT-Jojonomic-Indonesia/microkit/response"
)

type RetrieveInput struct {
	AggregateID   string `json:"aggregate_id"`
	AggregateName string `json:"aggregate_name"`
	AfterVersion  int    `json:"after_version"`
}

func RetrieveHandler(w http.ResponseWriter, r *http.Request) {
	dataInput := RetrieveInput{}
	err := json.NewDecoder(r.Body).Decode(&dataInput)
	if err != nil {
		response.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	events, err := esam.Retrieve(dataInput.AggregateID, dataInput.AggregateName, dataInput.AfterVersion)
	if err != nil {
		response.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	dataResp := map[string]any{
		"message": "success",
		"data":    events,
		"error":   false,
	}
	response.JSON(w, dataResp, http.StatusOK)

}
