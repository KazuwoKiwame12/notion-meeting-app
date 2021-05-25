package infrastructure

import (
	"app/config"
	"app/domain/adapter"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type SqlHandler struct {
	DB *sql.DB
}

func NewSqlHandler() (*SqlHandler, error) {
	db, err := sql.Open("postgres", config.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open PostgreSQL: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect PostgreSQL: %w", err)
	}

	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(100 * time.Second)

	return &SqlHandler{DB: db}, nil
}

func (sh *SqlHandler) Execute(statement string, args ...interface{}) (adapter.Result, error) {
	result, err := sh.DB.Exec(statement, args...)
	return result, err
}

func (sh *SqlHandler) Query(query string, args ...interface{}) (adapter.Rows, error) {
	rows, err := sh.DB.Query(query, args...)
	return rows, err
}

func (sh *SqlHandler) QueryRow(query string, args ...interface{}) adapter.Row {
	return sh.DB.QueryRow(query, args...)
}
