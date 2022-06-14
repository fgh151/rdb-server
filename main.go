package main

import (
	"context"
	"db-server/drivers"
	err2 "db-server/err"
	"db-server/migrations"
	"db-server/models"
	"db-server/server"
	"db-server/web"
	"flag"
	"github.com/evalphobia/logrus_sentry"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {

	verboseMode := flag.Bool("v", false, "Verbose mode")
	migrateFlag := flag.Bool("m", false, "Run migrations")
	demoFlag := flag.Bool("demo", false, "Fill demo data")
	docsFlag := flag.Bool("docs", true, "Disable public docs")
	sentryFlag := flag.Bool("sentry", true, "Disable sentry logs")

	flag.Parse()

	if *verboseMode {
		log.SetLevel(log.DebugLevel)
	}

	log.Debug("Init log system done")
	log.Debug("Init sentry")

	err := godotenv.Load()
	err2.PanicErr(err)

	if *sentryFlag {
		hook, err := logrus_sentry.NewWithTagsSentryHook(
			os.Getenv("SENTRY_DSN"),
			map[string]string{"ENVIRONMENT": os.Getenv("SENTRY_ENVIRONMENT")},
			[]log.Level{
				log.PanicLevel,
				log.FatalLevel,
				log.ErrorLevel,
			})
		if err == nil {
			log.AddHook(hook)
		}
	}

	log.Debug("Init meta db connection")
	db := server.MetaDb.GetConnection()

	if *migrateFlag {
		log.Debug("Try to migrate db")
		migrations.Migrate(db)
		os.Exit(0)
	}

	if *demoFlag {
		log.Debug("Fill db demo data")
		models.CreateDemo()
		os.Exit(0)
	}

	log.Debug("Init mongo db connection")

	client, _ := drivers.GetDbInstance().GetConnection()

	defer func() {
		log.Debug("Close mongo db connection")

		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	log.Debug("Start cron")
	models.InitCron()

	log.Debug("Init web server")
	web.InitServer(docsFlag)
}
