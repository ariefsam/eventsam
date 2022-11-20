package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func getRoutes() http.Handler {
	route := mux.NewRouter()

	route.Path("/health").Methods(http.MethodGet).HandlerFunc(HealthHandler)
	route.Path("/store").Methods(http.MethodPost).HandlerFunc(StoreHandler)
	route.Path("/retrieve").Methods(http.MethodPost).HandlerFunc(RetrieveHandler)
	return route
}
