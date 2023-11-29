package database

import (
	"fmt"

	"github.com/go-faster/errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Addr     string
	Port     uint16
	User     string
	Password string
	DB       string
}

func New(cfg *Config) (*sqlx.DB, error) {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.User, cfg.Password, cfg.Addr, cfg.Port, cfg.DB)
	conn, err := sqlx.Open("mysql", dataSource)
	if err != nil {
		return nil, errors.Wrap(err, "mysql-connection")
	}

	err = conn.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "mysql-ping-failed")
	}

	return conn, nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
