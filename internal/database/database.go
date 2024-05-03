package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DatabaseTimeoutSeconds int    = 5
	MigrationTableName     string = "_gomigrations"
	MigrationSchema        string = "CREATE TABLE IF NOT EXISTS %s (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);"
	MigrationsDir          string = "migrations"
)

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

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
			MigrationSchema,
			MigrationTableName,
		),
	)

	return err
}

func CheckTableExistence(db *sql.DB) (bool, error) {
	var tableCount int
	err := db.QueryRow(
		fmt.Sprintf(
			"SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='%s'",
			MigrationTableName,
		),
	).Scan(&tableCount)

	return tableCount > 0, err
}

func ApplyMigrations(db *sql.DB) error {
	files, err := os.ReadDir(MigrationsDir)
	if err != nil {
		return err
	}

	exists, err := CheckTableExistence(db)
	if err != nil {
		return err
	}

	if !exists {
		err = CreateMigrationTable(db)
		if err != nil {
			return err
		}
	}

	var appliedMigrations []string
	rows, err := db.Query(fmt.Sprintf("SELECT name FROM %s", MigrationTableName))
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return err
		}
		appliedMigrations = append(appliedMigrations, name)
	}

	for _, file := range files {
		if strings.Count(file.Name(), "down") > 0 {
			continue
		}
		if !contains(appliedMigrations, file.Name()) {
			path := filepath.Join(MigrationsDir, file.Name())
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			_, err = db.Exec(string(data))
			if err != nil {
				return err
			}

			tx, _ := db.Begin()

			_, err = tx.Exec(
				fmt.Sprintf("INSERT INTO %s (name) VALUES (?)", MigrationTableName),
				file.Name(),
			)
			if err != nil {
				return err
			}

			_, err = tx.Exec(string(data))
			if err != nil {
				return err
			}

			fmt.Printf("Applied migration: %s\n", path)

			tx.Commit()
		}
	}

	return nil
}

func rollbackMigration(db *sql.DB, migrationName string) error {
	path := filepath.Join(MigrationsDir, migrationName)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(data))
	if err != nil {
		return err
	}

	_, err = db.Exec(
		fmt.Sprintf("DELETE FROM %s WHERE name = ?", MigrationTableName),
		migrationName,
	)
	if err != nil {
		return err
	}

	fmt.Printf("Rolled back migration: %s\n", migrationName)
	return nil
}

func RollbackMigrations(db *sql.DB) error {
	var appliedMigrations []string
	rows, err := db.Query(fmt.Sprintf("SELECT name FROM %s ORDER BY id DESC", MigrationTableName))
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return err
		}
		appliedMigrations = append(appliedMigrations, strings.Replace(name, "up", "down", 1))
	}

	for _, migrationName := range appliedMigrations {
		err := rollbackMigration(db, migrationName)
		if err != nil {
			return err
		}
	}

	return nil
}
