package db

import (
	"db-server/drivers"
	"db-server/modules"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
)

type connection struct {
	db *gorm.DB
}

func (c connection) connect() *gorm.DB {
	log.Debug("Set meta db driver " + os.Getenv("META_DB_TYPE"))

	var db *gorm.DB
	var err error

	switch os.Getenv("META_DB_TYPE") {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(os.Getenv("META_DB_DSN")), &gorm.Config{})
	case "mysql":
		conn := drivers.NewMysqlConnectionFromEnv()
		db, err = gorm.Open(mysql.Open(conn.GetDsn()), &gorm.Config{})
	case "postgres":
		conn := drivers.NewPostgresConnectionFromEnv()
		db, err = gorm.Open(postgres.Open(conn.GetDsn()), &gorm.Config{})
	}

	if err != nil {
		panic("failed to connect database")
	}

	if log.GetLevel() >= log.DebugLevel {
		return db.Debug()
	}

	return db
}

func (c connection) GetConnection() *gorm.DB {

	if c.db == nil {
		c.db = c.connect()
	}

	return c.db
}

var MetaDb = connection{}

func (c connection) TotalRecords(m modules.Model) *int64 {
	conn := c.GetConnection()
	var cnt int64
	conn.Model(&m).Count(&cnt)
	return &cnt
}

func (c connection) ListQuery(limit int, offset int, sort string, order string, filter map[string]string, dest interface{}, preload []string) {

	query := c.GetConnection().Offset(offset).Limit(limit).Order(clause.OrderBy{Expression: clause.Expr{SQL: "? ?", Vars: []interface{}{[]string{sort, order}}}})

	if filter != nil && len(filter) > 0 {
		for k, v := range filter {
			query.Where(k+" = ?", v)
		}
	}

	for _, relation := range preload {
		query.Preload(relation)
	}

	query.Find(dest)
}
