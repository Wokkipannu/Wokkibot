package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: go run cmd/import/main.go <database_path> <names_file>")
	}

	databasepath := os.Args[1]
	filepath := os.Args[2]

	// Open database connection
	db, err := sql.Open("libsql", fmt.Sprintf("file:%s", databasepath))
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create names table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS names (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	)`)
	if err != nil {
		log.Fatalf("Failed to create names table: %v", err)
	}

	// Read and process the names file
	content, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Failed to read names file: %v", err)
	}

	names := strings.Split(strings.TrimSpace(string(content)), "\n")

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT OR IGNORE INTO names (name) VALUES (?)")
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}

		_, err = stmt.Exec(name)
		if err != nil {
			log.Fatalf("Failed to insert name %q: %v", name, err)
		}
	}

	if err = tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	fmt.Printf("Successfully imported %d names\n", len(names))
}
