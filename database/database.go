package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

var (
	db     *sql.DB
	dbOnce sync.Once
)

type Config struct {
	DatabaseURL string
}

func Initialize(cfg Config) error {
	var err error
	dbOnce.Do(func() {
		db, err = sql.Open("libsql", cfg.DatabaseURL)
		if err != nil {
			log.Printf("Failed to create database client: %v", err)
			return
		}

		if err = db.Ping(); err != nil {
			log.Printf("Failed to ping database: %v", err)
			return
		}

		err = initializeSchema()
		if err != nil {
			log.Printf("Failed to initialize schema: %v", err)
			return
		}

		err = runMigrations()
		if err != nil {
			log.Printf("Failed to run migrations: %v", err)
			return
		}
	})
	return err
}

func GetDB() *sql.DB {
	return db
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func initializeSchema() error {
	schemas := []string{
		`CREATE TABLE IF NOT EXISTS migrations (
    	version INTEGER PRIMARY KEY,
    	applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS custom_commands (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			guild_id TEXT NOT NULL,
			prefix TEXT NOT NULL,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			output TEXT NOT NULL,
			author TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS friday_clips (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS guilds (
			id TEXT PRIMARY KEY,
			trivia_token TEXT,
			pin_channel TEXT
		)`,
	}

	for _, schema := range schemas {
		_, err := db.Exec(schema)
		if err != nil {
			return fmt.Errorf("failed to create schema: %v", err)
		}
	}

	return nil
}
