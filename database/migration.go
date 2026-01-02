package database

import (
	"fmt"
	"log/slog"
	"strings"
)

type Migration struct {
	Version     int
	Description string
	SQL         string
	IgnoreError string
}

// TODO: move migrations to their own individual files perhaps. For now, this works, but in the long run this can get really big
var migrations = []Migration{
	{
		Version:     1,
		Description: "Add pin_channel to guilds",
		SQL:         "ALTER TABLE guilds ADD COLUMN pin_channel TEXT;",
		IgnoreError: "duplicate column",
	},
	{
		Version:     2,
		Description: "Add data to pizza_toppings table",
		SQL:         "INSERT INTO pizza_toppings (name) VALUES ('Ananas'), ('Aurajuusto'), ('Chili'), ('Jalopeno'), ('Tuplajuusto'), ('Kananmuna'), ('Katkarapu'), ('Kermaperunat'), ('Oliivi'), ('Pekoni'), ('Pippurikastike'), ('Punasipuli'), ('Salaatti'), ('Simpukka'), ('Smetana'), ('Tomaatti'), ('Herkkusieni'), ('Anjovis'), ('BBQ-kastike'), ('Fetajuusto'), ('Jauheliha'), ('Kana'), ('Kapris'), ('Kebab'), ('Mozzarella'), ('Paprika'), ('Pepperoni'), ('Pizzasuikale'), ('Rucola'), ('Salami'), ('Sipuli'), ('Suolakurkku'), ('Tonnikala'), ('Banaani'), ('Currykastike');",
	},
	{
		Version:     3,
		Description: "Add convert_x_links to guilds",
		SQL:         "ALTER TABLE guilds ADD COLUMN convert_x_links BOOLEAN DEFAULT TRUE;",
		IgnoreError: "duplicate column",
	},
	{
		Version:     4,
		Description: "Add reminders table",
		SQL:         "CREATE TABLE IF NOT EXISTS reminders (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id TEXT NOT NULL, channel_id TEXT NOT NULL, message TEXT NOT NULL, remind_at DATETIME NOT NULL);",
	},
	{
		Version:     5,
		Description: "Add guild_id to reminders",
		SQL:         "ALTER TABLE reminders ADD COLUMN guild_id TEXT;",
		IgnoreError: "duplicate column",
	},
	{
		Version:     6,
		Description: "Initialize statistics table with a single row",
		SQL:         "INSERT INTO statistics (video_downloads, names_given, songs_played, pizzas_generated, coins_flipped, dice_rolled, trivia_games_played, trivia_games_won, trivia_games_lost) SELECT 0, 0, 0, 0, 0, 0, 0, 0, 0 WHERE NOT EXISTS (SELECT 1 FROM statistics);",
	},
	{
		Version:     7,
		Description: "Add blackjack_games_played to statistics",
		SQL:         "ALTER TABLE statistics ADD COLUMN blackjack_games_played INTEGER DEFAULT 0;",
		IgnoreError: "duplicate column",
	},
}

func runMigrations() error {
	db := GetDB()

	var lastVersion int
	err := db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM migrations").Scan(&lastVersion)
	if err != nil {
		return fmt.Errorf("failed to get last migration version: %v", err)
	}

	for _, migration := range migrations {
		if migration.Version <= lastVersion {
			continue
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to start transaction: %v", err)
		}

		_, err = tx.Exec(migration.SQL)
		if err != nil {
			if migration.IgnoreError != "" && strings.Contains(err.Error(), migration.IgnoreError) {
				slog.Info("Migration already applied (ignoring expected error)",
					"version", migration.Version,
					"description", migration.Description)
				tx.Rollback()
				_, recordErr := db.Exec("INSERT INTO migrations (version) VALUES (?)", migration.Version)
				if recordErr != nil {
					return fmt.Errorf("failed to record migration %d: %v", migration.Version, recordErr)
				}
				continue
			}
			tx.Rollback()
			return fmt.Errorf("failed to apply migration %d: %v", migration.Version, err)
		}

		_, err = tx.Exec("INSERT INTO migrations (version) VALUES (?)", migration.Version)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %d: %v", migration.Version, err)
		}

		if err = tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d: %v", migration.Version, err)
		}

		slog.Info("Applied migration",
			"version", migration.Version,
			"description", migration.Description)
	}

	return nil
}
