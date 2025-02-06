package handlers

import (
	"fmt"
	"log/slog"
	"wokkibot/database"
	"wokkibot/types"

	"github.com/disgoorg/snowflake/v2"
)

func LoadCommands() ([]types.Command, error) {
	var commands []types.Command

	db := database.GetDB()

	rows, err := db.Query("SELECT name, prefix, description, output, author, guild_id FROM custom_commands")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var command types.Command
		err := rows.Scan(&command.Name, &command.Prefix, &command.Description, &command.Output, &command.Author, &command.GuildID)
		if err != nil {
			return nil, err
		}
		commands = append(commands, command)
	}

	slog.Info("Loaded custom commands", "count", len(commands))

	return commands, nil
}

func (h *Handler) AddOrUpdateCommand(newCommand types.Command) error {
	db := database.GetDB()

	var exists bool
	var existingAuthor snowflake.ID
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM custom_commands WHERE name = $1 AND prefix = $2)", newCommand.Name, newCommand.Prefix).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		err = db.QueryRow("SELECT author FROM custom_commands WHERE name = $1 AND prefix = $2", newCommand.Name, newCommand.Prefix).Scan(&existingAuthor)
		if err != nil {
			return err
		}

		if existingAuthor != newCommand.Author {
			return fmt.Errorf("only the original author can update this command")
		}

		_, err = db.Exec("UPDATE custom_commands SET description = $1, output = $2, guild_id = $3 WHERE name = $4 AND prefix = $5",
			newCommand.Description, newCommand.Output, newCommand.GuildID, newCommand.Name, newCommand.Prefix)

		for i, cmd := range h.CustomCommands {
			if cmd.Name == newCommand.Name && cmd.Prefix == newCommand.Prefix {
				h.CustomCommands[i] = newCommand
				break
			}
		}
	} else {
		_, err = db.Exec("INSERT INTO custom_commands (name, prefix, description, output, author, guild_id) VALUES ($1, $2, $3, $4, $5, $6)",
			newCommand.Name, newCommand.Prefix, newCommand.Description, newCommand.Output, newCommand.Author, newCommand.GuildID)

		h.CustomCommands = append(h.CustomCommands, newCommand)
	}

	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) RemoveCommand(prefix string, name string, author snowflake.ID) error {
	db := database.GetDB()

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM custom_commands WHERE prefix = $1 AND name = $2 AND author = $3)", prefix, name, author).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("command does not exist or you are not the author")
	}

	_, err = db.Exec("DELETE FROM custom_commands WHERE prefix = $1 AND name = $2 AND author = $3", prefix, name, author)
	if err != nil {
		return err
	}

	for i, cmd := range h.CustomCommands {
		if cmd.Name == name && cmd.Prefix == prefix && cmd.Author == author {
			h.CustomCommands = append(h.CustomCommands[:i], h.CustomCommands[i+1:]...)
			break
		}
	}

	return nil
}
