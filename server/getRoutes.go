package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func getRoutes() http.Handler {
	route := mux.NewRouter()

	route.Path("/health").Methods(http.MethodGet).HandlerFunc(HealthHandler)
	route.Path("/store").Methods(http.MethodPost).HandlerFunc(StoreHandler)
	route.Path("/retrieve").Methods(http.MethodPost).HandlerFunc(RetrieveHandler)
	route.Path("/fetch-all-event").Methods(http.MethodPost).HandlerFunc(FetchAllEventHandler)
	route.Path("/fetch-aggregate-event").Methods(http.MethodPost).HandlerFunc(FetchAggregateEventHandler)
	return route
}
