package main

import (
	"context"
	"db-server/drivers"
	err2 "db-server/err"
	"db-server/meta"
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

	client, _ := drivers.GetDbInstance().GetConnection()
	meta.MetaDb.GetConnection()

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	web.InitServer()
}
