package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DatabaseTimeoutSeconds int    = 5
	MigrationTableName     string = "_gomigrations"
)

func NewDatabaseConn(db_url string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite3", db_url)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func CreateMigrationTable(db *sql.DB) error {
	migrationCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(
		migrationCtx,
		fmt.Sprintf(
			"CREATE TABLE IF NOT EXISTS %s (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);",
			MigrationTableName,
		),
	)

	return err
}
