package server

import (
	"log"
	"os"
	"sync"

	"github.com/PT-Jojonomic-Indonesia/microkit/server"
	"github.com/ariefsam/eventsam"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var esam eventsam.Eventsam
var isSlave bool
var logger = log.New(os.Stdout, "Eventsam v.1.0.0 ", log.LstdFlags|log.Llongfile)

func Serve(db *gorm.DB) {
	var err error

	log.SetFlags(log.LstdFlags | log.Llongfile)
	godotenv.Load()

	cond = sync.NewCond(&sync.Mutex{})

	if os.Getenv("MASTER_ADDRESS") != "" {
		isSlave = true
		go SlaveSync(db)
	}

	esam, err = eventsam.NewEventsam(db)
	if err != nil {
		log.Println(err)
		return
	}

	router := getRoutes()
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	logger.Println("Server is running on port " + port)
	server.Serve(port, router)
}
