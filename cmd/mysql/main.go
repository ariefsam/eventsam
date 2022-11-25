package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ariefsam/eventsam/server"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	godotenv.Load()
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	database := os.Getenv("MYSQL_DATABASE")
	dbConnection := fmt.Sprintf("%[4]s:%[5]s@tcp(%[1]s:%[3]s)/%[2]s?charset=utf8mb4&parseTime=True&loc=Local", host, database, port, user, password)

	logService := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,         // Disable color
		})
	db, err := gorm.Open(mysql.Open(dbConnection), &gorm.Config{
		Logger: logService,
	})
	if err != nil {
		log.Println("Failed connect to Database", err)
		return
	}
	if err != nil {
		log.Fatal(err)
		return
	}
	server.Serve(db)
}
