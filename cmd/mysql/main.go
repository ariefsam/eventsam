package main

import (
	"log"
	"os"
	"time"

	"github.com/ariefsam/eventsam/server"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
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
	if err != nil {
		log.Fatal(err)
		return
	}
	server.Serve(db)
}
