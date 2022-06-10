package drivers

import (
	"fmt"
	"os"
)

// MysqlConnection Mysql connection description
type MysqlConnection struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
}

// NewMysqlConnectionFromEnv Create Mysql connection from env
func NewMysqlConnectionFromEnv() MysqlConnection {
	return MysqlConnection{
		Host:     os.Getenv("META_DB_HOST"),
		Port:     os.Getenv("META_DB_PORT"),
		User:     os.Getenv("META_DB_USER"),
		Password: os.Getenv("META_DB_PASSWORD"),
		DbName:   os.Getenv("META_DB_DBNAME"),
	}
}

// GetDsn Get mysql connection DSN
func (c MysqlConnection) GetDsn() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DbName,
	)
}
