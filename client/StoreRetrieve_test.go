package client_test

import (
	"log"
	"testing"

	"github.com/ariefsam/eventsam/client"
	"github.com/stretchr/testify/assert"
)

func TestNewEventsam(t *testing.T) {
	es, err := client.NewEventsam("http://localhost:8009")
	assert.NoError(t, err)
	data := map[string]any{
		"a": 10,
		"b": "c",
	}
	resp, err := es.Store("aggregate_id", "aggregate_type", "event_type", 2, data)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	log.Println(resp)

	resp2, err := es.Retrieve("aggregate_id", "aggregate_type", 0)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp2)
	log.Println(resp2)
}
