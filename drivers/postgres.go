package drivers

import (
	err2 "db-server/err"
	"fmt"
	"os"
	"strconv"
)

type PostgresConnection struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

func NewPostgresConnectionFromEnv() PostgresConnection {

	port, err := strconv.Atoi(os.Getenv("META_DB_PORT"))
	if err != nil {
		err2.PanicErr(err)
	}

	return PostgresConnection{
		Host:     os.Getenv("META_DB_HOST"),
		Port:     port,
		User:     os.Getenv("META_DB_USER"),
		Password: os.Getenv("META_DB_PASSWORD"),
		DbName:   os.Getenv("META_DB_DBNAME"),
	}
}

func (c PostgresConnection) GetDsn() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		c.Host,
		c.User,
		c.Password,
		c.DbName,
		c.Port,
	)
}
