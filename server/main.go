package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/PT-Jojonomic-Indonesia/microkit/server"
	"github.com/ariefsam/eventsam"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var esam eventsam.Eventsam
var isSlave bool
var messages chan string
var cond *sync.Cond

func main() {

	log.SetFlags(log.LstdFlags | log.Llongfile)
	godotenv.Load()

	cond = sync.NewCond(&sync.Mutex{})

	filepath := os.Getenv("DB_FILEPATH")
	if filepath == "" {
		filepath = "./event.db"
	}

	logService := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,         // Disable color
		})
	db, err := gorm.Open(sqlite.Open(filepath), &gorm.Config{
		Logger: logService,
	})
	if err != nil {
		log.Fatal(err)
		return
	}
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
	server.Serve(port, router)
}
