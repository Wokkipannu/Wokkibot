package wokkibot

import (
	"encoding/json"
	"os"

	"github.com/disgoorg/snowflake/v2"
)

type Command struct {
	Name        string       `json:"name"`
	Prefix      string       `json:"prefix"`
	Description string       `json:"description"`
	Output      string       `json:"output"`
	Author      snowflake.ID `json:"author"`
}

func LoadCommands(filename string) ([]Command, error) {
	var commands []Command
	file, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return commands, nil
		}
		return nil, err
	}
	err = json.Unmarshal(file, &commands)
	if err != nil {
		return nil, err
	}
	return commands, nil
}

func SaveCommands(filename string, commands []Command) error {
	file, err := json.MarshalIndent(commands, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, file, 0644)
	if err != nil {
		return err
	}
	return nil
}

func AddCommand(filename string, newCommand Command) error {
	commands, err := LoadCommands(filename)
	if err != nil {
		return err
	}
	commands = append(commands, newCommand)
	return SaveCommands(filename, commands)
}

func AddOrUpdateCommand(filename string, newCommand Command) error {
	commands, err := LoadCommands(filename)
	if err != nil {
		return err
	}
	for i, cmd := range commands {
		if cmd.Prefix == newCommand.Prefix && cmd.Name == newCommand.Name {
			commands[i] = newCommand
			return SaveCommands(filename, commands)
		}
	}

	commands = append(commands, newCommand)
	return SaveCommands(filename, commands)
}

func RemoveCommand(filename string, prefix string, name string) error {
	commands, err := LoadCommands(filename)
	if err != nil {
		return err
	}

	updatedCommands := commands[:0]
	for _, cmd := range commands {
		if !(cmd.Prefix == prefix && cmd.Name == name) {
			updatedCommands = append(updatedCommands, cmd)
		}
	}

	return SaveCommands(filename, updatedCommands)
}
