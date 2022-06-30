package ds

import (
	err2 "db-server/err"
	"db-server/modules/project"
	"db-server/server/db"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type DsType string

const (
	DSTypeMysql    DsType = "Mysql"
	DSTypePostgres DsType = "Postgres"
	DSTypeSqlite   DsType = "Sqlite"

	DSTypeXML DsType = "Xml"
)

type DataSource struct {
	// The data source UUID
	// example: 6204011c-30e6-408b-8aaa-dd8219860b4b
	Id uuid.UUID `gorm:"primarykey" json:"id"`
	// The data source type
	// Enum of DsType
	// example: Mysql
	Type DsType `json:"type"`
	// Data source title
	Title string `json:"title"`
	// Data source dsn
	Dsn string `json:"dsn"`
	// Linked project  UUID
	// example: 6204011c-30e6-408b-8aaa-dd8214860b4b
	ProjectId uuid.UUID `json:"project_id"`
	// Cache result in local db
	Cache bool `json:"cache"`
	// Linked project
	Project   project.Project
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Endpoints []DataSourceEndpoint `json:"endpoints"`
}

// TableName Gorm table name
func (p DataSource) TableName() string {
	return "ds_source"
}

func (p DataSource) List(limit int, offset int, sort string, order string, filter map[string]string) []interface{} {
	var sources []DataSource

	db.MetaDb.ListQuery(limit, offset, sort, order, filter, &sources, make([]string, 0))

	y := make([]interface{}, len(sources))
	for i, v := range sources {
		y[i] = v
	}

	return y
}

func (p DataSource) Total() *int64 {
	return db.MetaDb.TotalRecords(&DataSource{})
}

func (p DataSource) GetById(id string) interface{} {
	var source DataSource

	conn := db.MetaDb.GetConnection()

	conn.First(&source, "id = ?", id)

	return source
}

func (p DataSource) Delete(id string) {
	conn := db.MetaDb.GetConnection()
	conn.Where("id = ?", id).Delete(&p)
}

type DataSourceEndpoint struct {
	// The data source endpoint UUID
	// example: 6234011c-30e6-408b-8aaa-dd8219860b4b
	Id uuid.UUID `gorm:"primarykey" json:"id"`
	// Data source endpoint title
	Title string `json:"title"`
	// Endpoint table name
	SourceTable string `json:"table_name"`
	// Linked data source UUID
	// example: 6204011c-33e6-408b-8aaa-dd8214860b4b
	DataSourceId uuid.UUID `json:"data_source"`
	// Linked data source
	DataSource DataSource
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName Gorm table name
func (e DataSourceEndpoint) TableName() string {
	return "ds_endpoint"
}

func (e DataSourceEndpoint) List(limit int, offset int, sort string, order string, filter map[string]string) []interface{} {
	var sources []DataSourceEndpoint

	db.MetaDb.ListQuery(limit, offset, sort, order, filter, &sources, make([]string, 0))

	y := make([]interface{}, len(sources))
	for i, v := range sources {
		y[i] = v
	}

	return y
}

var dsConnections = make(map[string]*gorm.DB)

func (e DataSourceEndpoint) getConnection() (*gorm.DB, error) {
	if conn, ok := dsConnections[e.DataSource.Id.String()]; ok {
		return conn, nil
	}

	switch e.DataSource.Type {
	case DSTypeMysql:
		conn, err := gorm.Open(mysql.Open(e.DataSource.Dsn), &gorm.Config{})
		return e.attachConnectionToPool(conn, err)

	case DSTypePostgres:
		conn, err := gorm.Open(postgres.Open(e.DataSource.Dsn), &gorm.Config{})
		return e.attachConnectionToPool(conn, err)

	case DSTypeSqlite:
		conn, err := gorm.Open(sqlite.Open(e.DataSource.Dsn), &gorm.Config{})
		return e.attachConnectionToPool(conn, err)

	case DSTypeXML:
		return nil, nil
	}

	return nil, nil
}

func (e DataSourceEndpoint) attachConnectionToPool(conn *gorm.DB, err error) (*gorm.DB, error) {
	if err != nil {
		return nil, err
	}
	dsConnections[e.DataSource.Id.String()] = conn
	return conn, nil
}

func (e DataSourceEndpoint) GetById(id string) interface{} {
	var source DataSourceEndpoint

	conn := db.MetaDb.GetConnection()

	conn.First(&source, "id = ?", id)

	return source
}

func (e DataSourceEndpoint) Total() *int64 {

	conn, err := e.getConnection()

	var cnt int64 = 0

	if err != nil {
		err2.DebugErr(err)
		return &cnt
	}

	conn.Table(e.SourceTable).Count(&cnt)

	return &cnt
}

func (e DataSourceEndpoint) Delete(id string) {
	conn := db.MetaDb.GetConnection()
	conn.Where("id = ?", id).Delete(&e)
}
