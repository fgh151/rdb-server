package main

import (
	"db-server/db"
	err2 "db-server/err"
	"db-server/models"
	"db-server/web"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
)

func main() {
	err := godotenv.Load()
	err2.CheckErr(err)

	err = sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	dbInstance := db.DB.GetConnection()

	err = dbInstance.AutoMigrate(&models.Project{})
	err2.CheckErr(err)
	err = dbInstance.AutoMigrate(&models.User{})
	err2.CheckErr(err)
	err = dbInstance.AutoMigrate(&models.Topic{})
	err2.CheckErr(err)
	err = dbInstance.AutoMigrate(&models.Message{})
	err2.CheckErr(err)

	web.InitServer()
}
