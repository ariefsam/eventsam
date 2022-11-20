package main

import (
	"net/http"

	"github.com/PT-Jojonomic-Indonesia/microkit/response"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	resData := map[string]interface{}{
		"status": "ok",
	}
	response.JSON(w, resData, http.StatusOK)
}
