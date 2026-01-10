package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/newrelic/go-agent/v3/integrations/nrmysql"
)

func Connect(dsn string) (*sql.DB, error) {
	// Используем nrmysql драйвер для автоматического трейсинга SQL-запросов в New Relic
	db, err := sql.Open("nrmysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
