package main

import (
	"context"
	"db-server/drivers"
	err2 "db-server/err"
	"db-server/meta"
	"db-server/migrations"
	"db-server/models"
	"db-server/web"
	"flag"
	"github.com/joho/godotenv"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
)

func main() {

	verboseMode := flag.Bool("v", false, "Verbose mode")
	migrateFlag := flag.Bool("m", false, "Run migrations")
	demoFlag := flag.Bool("demo", false, "Fill demo data")
	flag.Parse()

	if *verboseMode {
		log.SetLevel(log.DebugLevel)
	}

	log.Debug("Init log system done")
	log.Debug("Init sentry")

	err := godotenv.Load()
	err2.PanicErr(err)

	err = sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	log.Debug("Init mongo db connection")

	client, _ := drivers.GetDbInstance().GetConnection()

	log.Debug("Init meta db connection")
	db := meta.MetaDb.GetConnection()

	if *migrateFlag {
		log.Debug("Try to migrate db")
		migrations.Migrate(db)
	}

	if *demoFlag {
		log.Debug("Fill db demo data")
		models.CreateDemo()
	}

	defer func() {
		log.Debug("Close mongo db connection")

		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	log.Debug("Init web server")
	web.InitServer()
}
