package server

import (
	"github.com/robfig/cron/v3"
)

type ServerCron struct {
	Cron *cron.Cron
}

var Cron = ServerCron{}

func (sc ServerCron) GetScheduler() *cron.Cron {
	if sc.Cron == nil {
		sc.Cron = cron.New()
	}

	return sc.Cron
}
