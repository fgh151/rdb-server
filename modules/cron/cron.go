package cron

import (
	"db-server/modules/cf"
	"db-server/server"
	"db-server/server/db"
	"db-server/utils"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type CronJob struct {
	Id         uuid.UUID        `gorm:"primarykey" json:"id"`
	Title      string           `json:"title"`
	TimeParams string           `json:"time_params"`
	FunctionId uuid.UUID        `json:"function_id"`
	Function   cf.CloudFunction `gorm:"foreignKey:FunctionId" json:"function"`
	CronId     cron.EntryID     `gorm:"index" json:"-"`
	CreatedAt  time.Time        `json:"-"`
	UpdatedAt  time.Time        `json:"-"`
	DeletedAt  gorm.DeletedAt   `gorm:"index" json:"-"`
}

// TableName Gorm table name
func (j CronJob) TableName() string {
	return "cron_job"
}

func (j CronJob) List(limit int, offset int, sort string, order string, filter map[string]string) ([]interface{}, error) {
	var jobs []CronJob

	db.MetaDb.ListQuery(limit, offset, sort, order, filter, &jobs, make([]string, 0))

	y := make([]interface{}, len(jobs))
	for i, v := range jobs {
		y[i] = v
	}

	return y, nil
}

func (j CronJob) GetById(id string) (interface{}, error) {
	var job CronJob

	conn := db.MetaDb.GetConnection()

	conn.First(&job, "id = ?", id)

	return job, nil
}

func (j CronJob) Delete(id string) {
	conn := db.MetaDb.GetConnection()
	conn.Where("id = ?", id).Delete(&j)
}

func (j CronJob) Total() *int64 {
	return db.MetaDb.TotalRecords(&CronJob{})
}

func (j CronJob) Schedule(cron *cron.Cron) {
	var err error

	j.CronId, err = cron.AddFunc(j.TimeParams, func() {
		log.Debug("Run cron " + utils.CleanInputString(j.Id.String()))
		function, err := cf.CloudFunction{}.GetById(j.FunctionId.String())
		if err == nil {
			id, _ := uuid.NewUUID()
			function.(cf.CloudFunction).Run(id)
		}
	})

	if err != nil {
		log.Debug(err)
	} else {
		db.MetaDb.GetConnection().Save(&j)
	}
}

func InitCron() {
	log.Debug("Start cron")
	c := server.Cron.GetScheduler()
	c.Start()

	offset := 0
	batchSize := 20
	var jobs []interface{}

	for {
		jobs, _ = CronJob{}.List(batchSize, offset, "id", "ASC", nil)

		if len(jobs) <= 0 {
			break
		}

		for _, job := range jobs {
			log.Debug("Add cron Job " + job.(CronJob).Id.String())
			job.(CronJob).Schedule(c)
		}

		offset += batchSize
	}
}

func StopCron() {
	c := server.Cron.GetScheduler()
	log.Debug("Stop cron")
	c.Stop()
}
