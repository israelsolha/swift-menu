package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func NewDb(cfg Config) (*sql.DB, error) {
	connectionStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Schema)
	return sql.Open("mysql", connectionStr)
}
