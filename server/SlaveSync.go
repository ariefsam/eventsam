package server

import (
	"log"
	"os"
	"time"

	"github.com/ariefsam/eventsam"
	"github.com/ariefsam/eventsam/client"
	"gorm.io/gorm"
)

func SlaveSync(db *gorm.DB) {
	for {

		curEvent := eventsam.EventEntity{}
		db.Last(&curEvent)

		afterID := curEvent.ID

		clientService, err := client.NewEventsam(os.Getenv("MASTER_ADDRESS"))
		if err != nil {
			log.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}
		events, err := clientService.FetchAllEvent(afterID, 100)
		if err != nil {
			log.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}
		if len(events) == 0 {
			time.Sleep(1 * time.Second)
		}
		eventMap := map[string]bool{}
		for _, event := range events {
			eventMap[event.AggregateName] = true
			err = db.Save(&event).Error
			if err != nil {
				log.Println(err)
				continue
			}
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
			for aggregateName := range eventMap {
				condA := getCondAggregate(aggregateName)
				condA.L.Lock()
				condA.Broadcast()
				condA.L.Unlock()
			}
		}()

	}
}
