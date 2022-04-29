package models

import (
	err2 "db-server/err"
	"db-server/meta"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type DsType string

const (
	DSTypeMysql    DsType = "Mysql"
	DSTypePostgres DsType = "Postgres"
)

type DataSource struct {
	Id        uuid.UUID `gorm:"primarykey" json:"id"`
	Type      DsType    `json:"type"`
	Title     string    `json:"title"`
	Dsn       string    `json:"dsn"`
	ProjectId uuid.UUID `json:"project_id"`
	Cache     bool      `json:"cache"`
	Project   Project
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p DataSource) List(limit int, offset int, sort string, order string) []interface{} {
	var sources []DataSource

	conn := meta.MetaDb.GetConnection()

	conn.Find(&sources).Limit(limit).Offset(offset).Order(order + " " + sort)

	y := make([]interface{}, len(sources))
	for i, v := range sources {
		y[i] = v
	}

	return y
}

func (p DataSource) Total() *int64 {
	conn := meta.MetaDb.GetConnection()
	var sources []DataSource
	var cnt int64
	conn.Find(&sources).Count(&cnt)

	return &cnt
}

func (p DataSource) GetById(id string) interface{} {
	var source DataSource

	conn := meta.MetaDb.GetConnection()

	conn.First(&source, "id = ?", id)

	return source
}

func (p DataSource) Delete(id string) {
	conn := meta.MetaDb.GetConnection()
	conn.Where("id = ?", id).Delete(&p)
}

type DataSourceEndpoint struct {
	Id    uuid.UUID `gorm:"primarykey" json:"id"`
	Title string    `json:"title"`

	TableName string `json:"table_name"`

	DataSourceId uuid.UUID `json:"data_source"`
	DataSource   DataSource
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (e DataSourceEndpoint) List(limit int, offset int, order string, sort string) []interface{} {
	arr := []interface{}{}

	conn, err := e.getConnection()

	if err != nil {
		err2.DebugErr(err)
		return arr
	}

	rows, err := conn.Debug().Table(e.TableName).Limit(limit).Offset(offset).Order(order + " " + sort).Rows()
	err2.DebugErr(err)

	defer func() { _ = rows.Close() }()

	for rows.Next() {

		cols, _ := rows.Columns()
		data := make(map[string]string)

		columns := make([]string, len(cols))
		columnPointers := make([]interface{}, len(cols))

		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		err := rows.Scan(columnPointers...)
		err2.DebugErr(err)

		for i, colName := range cols {
			data[colName] = columns[i]
		}

		arr = append(arr, data)
	}

	return arr
}

var dsConnections = make(map[string]*gorm.DB)

func (e DataSourceEndpoint) getConnection() (*gorm.DB, error) {
	if conn, ok := dsConnections[e.DataSource.Id.String()]; ok {
		return conn, nil
	}

	switch e.DataSource.Type {
	case DSTypeMysql:
		conn, err := gorm.Open(mysql.Open(e.DataSource.Dsn), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		dsConnections[e.DataSource.Id.String()] = conn
		return conn, nil
	}

	return nil, nil
}

func (e DataSourceEndpoint) GetById(id string) interface{} {
	var source DataSourceEndpoint

	conn := meta.MetaDb.GetConnection()

	conn.Preload("DataSource").First(&source, "id = ?", id)

	return source
}

func (e DataSourceEndpoint) Total() *int64 {

	conn, err := e.getConnection()

	var cnt int64 = 0

	if err != nil {
		err2.DebugErr(err)
		return &cnt
	}

	conn.Table(e.TableName).Count(&cnt)

	return &cnt
}
