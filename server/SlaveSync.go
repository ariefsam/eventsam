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
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in SlaveSync", r)
		}
		log.Println("Exiting")
	}()
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
		events, err := clientService.FetchAllEvent(int(afterID), 100)
		if err != nil {
			log.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}
		if len(events) == 0 {
			time.Sleep(1 * time.Second)
		}
		for _, event := range events {
			err = db.Save(&event).Error
			if err != nil {
				log.Println(err)
				continue
			}
		}
		go func() {
			cond.L.Lock()
			cond.Broadcast()
			cond.L.Unlock()
		}()

	}
}
