package models

import (
	"db-server/server"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type CronJob struct {
	Id         uuid.UUID      `gorm:"primarykey" json:"id"`
	Title      string         `json:"title"`
	TimeParams string         `json:"time_params"`
	FunctionId uuid.UUID      `json:"function_id"`
	Function   CloudFunction  `gorm:"foreignKey:FunctionId" json:"function"`
	CronId     cron.EntryID   `gorm:"index" json:"-"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (j CronJob) List(limit int, offset int, sort string, order string, filter map[string]interface{}) []interface{} {
	var jobs []CronJob

	conn := server.MetaDb.GetConnection()

	conn.Limit(limit).Offset(offset).Order(sort + " " + order).Where(filter).Find(&jobs)

	y := make([]interface{}, len(jobs))
	for i, v := range jobs {
		y[i] = v
	}

	return y
}

func (j CronJob) GetById(id string) interface{} {
	var job CronJob

	conn := server.MetaDb.GetConnection()

	conn.First(&job, "id = ?", id)

	return job
}

func (j CronJob) Delete(id string) {
	conn := server.MetaDb.GetConnection()
	conn.Where("id = ?", id).Delete(&j)
}

func (j CronJob) Total() *int64 {
	return TotalRecords(&CronJob{})
}

func (j CronJob) Schedule(cron *cron.Cron) {
	var err error

	j.CronId, err = cron.AddFunc(j.TimeParams, func() {
		log.Debug("Run cron " + j.Id.String())
		function := CloudFunction{}.GetById(j.FunctionId.String()).(CloudFunction)
		id, _ := uuid.NewUUID()
		function.Run(id)
	})

	if err != nil {
		log.Debug(err)
	} else {
		server.MetaDb.GetConnection().Save(&j)
	}
}

func InitCron() {
	c := server.Cron.GetScheduler()
	c.Start()
	defer func() {
		log.Debug("Stop cron")
		c.Stop()
	}()

	offset := 0
	batchSize := 20
	var jobs []interface{}

	for {
		jobs = CronJob{}.List(batchSize, offset, "id", "ASC", nil)

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
