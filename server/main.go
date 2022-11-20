package main

import (
	"eventsam"
	"log"
	"os"
	"time"

	"github.com/PT-Jojonomic-Indonesia/microkit/server"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var esam eventsam.Eventsam

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	godotenv.Load()

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

	// See "Important settings" section.

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
