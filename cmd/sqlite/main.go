package main

import (
	"log"
	"os"
	"time"

	"github.com/ariefsam/eventsam/server"
	"github.com/glebarez/sqlite"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	godotenv.Load()
	filepath := os.Getenv("DB_FILEPATH")
	if filepath == "" {
		filepath = "./event.db"
	}

	filepath += "?cache=shared&&mode=rwc"

	logService := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,         // Disable color
		})
	db, err := gorm.Open(sqlite.Open(filepath), &gorm.Config{
		Logger:      logService,
		PrepareStmt: true,
	})

	if err != nil {
		log.Fatal(err)
		return
	}
	sqldb, err := db.DB()
	if err != nil {
		log.Fatal(err)
		return
	}
	sqldb.SetMaxIdleConns(1)
	server.Serve(db)
}
