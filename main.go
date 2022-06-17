package main

import (
	"context"
	"db-server/drivers"
	err2 "db-server/err"
	"db-server/models"
	"db-server/server"
	"db-server/web"
	"flag"
	"github.com/evalphobia/logrus_sentry"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

// @title           Db server API
// @version         1.0
// @description     Db server API description.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    https://openitstudio.ru
// @contact.email  fedor@support-pc.org

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.bearer  BearerAuth

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @description					Description for what is this security definition being used
func main() {

	verboseMode := flag.Bool("v", false, "Verbose mode")
	docsFlag := flag.Bool("docs", true, "Disable public docs")
	sentryFlag := flag.Bool("sentry", true, "Disable sentry logs")
	mongoFlag := flag.Bool("mongo", true, "Disable mongo initialization")

	flag.Parse()

	if *verboseMode {
		log.SetLevel(log.DebugLevel)
	}

	log.Debug("Init log system done")

	err := godotenv.Load()
	err2.PanicErr(err)

	if *sentryFlag {
		log.Debug("Init sentry")
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
	server.MetaDb.GetConnection()

	if *mongoFlag {
		log.Debug("Init mongo db connection")

		client, _ := drivers.GetDbInstance().GetConnection()

		defer func() {
			log.Debug("Close mongo db connection")

			if err := client.Disconnect(context.TODO()); err != nil {
				panic(err)
			}
		}()
	}

	models.InitCron()

	log.Debug("Init web server")
	web.InitServer(docsFlag)
}
