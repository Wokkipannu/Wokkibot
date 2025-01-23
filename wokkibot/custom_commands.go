package wokkibot

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"
	"wokkibot/database"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"golang.org/x/exp/rand"
)

type Command struct {
	Name        string       `json:"name"`
	Prefix      string       `json:"prefix"`
	Description string       `json:"description"`
	Output      string       `json:"output"`
	Author      snowflake.ID `json:"author"`
	GuildID     snowflake.ID `json:"guild_id"`
}

func (b *Wokkibot) LoadCommands() ([]Command, error) {
	var commands []Command

	db := database.GetDB()

	rows, err := db.Query("SELECT name, prefix, description, output, author, guild_id FROM custom_commands")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var command Command
		err := rows.Scan(&command.Name, &command.Prefix, &command.Description, &command.Output, &command.Author, &command.GuildID)
		if err != nil {
			return nil, err
		}
		commands = append(commands, command)
	}

	slog.Info("Loaded custom commands", "count", len(commands))

	return commands, nil
}

func (b *Wokkibot) AddOrUpdateCommand(newCommand Command) error {
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

		for i, cmd := range b.CustomCommands {
			if cmd.Name == newCommand.Name && cmd.Prefix == newCommand.Prefix {
				b.CustomCommands[i] = newCommand
				break
			}
		}
	} else {
		_, err = db.Exec("INSERT INTO custom_commands (name, prefix, description, output, author, guild_id) VALUES ($1, $2, $3, $4, $5, $6)",
			newCommand.Name, newCommand.Prefix, newCommand.Description, newCommand.Output, newCommand.Author, newCommand.GuildID)

		b.CustomCommands = append(b.CustomCommands, newCommand)
	}

	if err != nil {
		return err
	}

	return nil
}

func (b *Wokkibot) RemoveCommand(prefix string, name string, author snowflake.ID) error {
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

	for i, cmd := range b.CustomCommands {
		if cmd.Name == name && cmd.Prefix == prefix && cmd.Author == author {
			b.CustomCommands = append(b.CustomCommands[:i], b.CustomCommands[i+1:]...)
			break
		}
	}

	return nil
}

func HandleCustomCommand(b *Wokkibot, e *events.MessageCreate) {
	input := e.Message.Content

	if input == "" {
		return
	}
	prefix := string(input[0])
	name := strings.TrimPrefix(input, prefix)

	for _, cmd := range b.CustomCommands {
		if cmd.Prefix == prefix && cmd.Name == name && cmd.GuildID == *e.GuildID {
			output := handleVariables(cmd.Output)

			e.Client().Rest().CreateMessage(e.Message.ChannelID, discord.NewMessageCreateBuilder().
				SetContent(output).
				SetMessageReferenceByID(e.Message.ID).
				SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false}).
				Build())
		}
	}
}

func handleVariables(text string) string {
	re := regexp.MustCompile(`\{\{(\w+)\|([^}]+)\}\}`)

	return re.ReplaceAllStringFunc(text, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match
		}

		variable := parts[1]
		value := parts[2]

		switch variable {
		case "time":
			loc, err := time.LoadLocation(value)
			if err != nil {
				slog.Error("Failed to load timezone", "location", value, "error", err)
				return "INVALID TIMEZONE NAME"
			}
			return time.Now().In(loc).Format("15:04 MST")
		case "random":
			choices := strings.Split(value, ";")
			if len(choices) == 0 {
				return "NO CHOICES PROVIDED"
			}
			randomIndex := rand.Intn(len(choices))
			return strings.TrimSpace(choices[randomIndex])
		default:
			return match
		}
	})
}
