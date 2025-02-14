package web

import (
	"database/sql"
	"wokkibot/database"
)

type CustomCommand struct {
	ID          int64  `json:"id"`
	GuildID     string `json:"guild_id"`
	Prefix      string `json:"prefix"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Output      string `json:"output"`
	Author      string `json:"author"`
}

type FridayClip struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

type PizzaTopping struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Custom Commands
func GetCustomCommands() ([]CustomCommand, error) {
	db := database.GetDB()

	rows, err := db.Query("SELECT id, guild_id, prefix, name, description, output, author FROM custom_commands")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []CustomCommand
	for rows.Next() {
		var cmd CustomCommand
		err := rows.Scan(&cmd.ID, &cmd.GuildID, &cmd.Prefix, &cmd.Name, &cmd.Description, &cmd.Output, &cmd.Author)
		if err != nil {
			return nil, err
		}
		commands = append(commands, cmd)
	}
	return commands, nil
}

func AddCustomCommand(cmd CustomCommand) error {
	db := database.GetDB()

	_, err := db.Exec(
		"INSERT INTO custom_commands (guild_id, prefix, name, description, output, author) VALUES (?, ?, ?, ?, ?, ?)",
		cmd.GuildID, cmd.Prefix, cmd.Name, cmd.Description, cmd.Output, cmd.Author,
	)
	return err
}

func UpdateCustomCommand(cmd CustomCommand) error {
	db := database.GetDB()

	result, err := db.Exec(
		"UPDATE custom_commands SET guild_id=?, prefix=?, name=?, description=?, output=?, author=? WHERE id=?",
		cmd.GuildID, cmd.Prefix, cmd.Name, cmd.Description, cmd.Output, cmd.Author, cmd.ID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func DeleteCustomCommand(id int64) error {
	db := database.GetDB()

	result, err := db.Exec("DELETE FROM custom_commands WHERE id=?", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Friday Clips
func GetFridayClips() ([]FridayClip, error) {
	db := database.GetDB()

	rows, err := db.Query("SELECT id, url FROM friday_clips")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clips []FridayClip
	for rows.Next() {
		var clip FridayClip
		err := rows.Scan(&clip.ID, &clip.URL)
		if err != nil {
			return nil, err
		}
		clips = append(clips, clip)
	}
	return clips, nil
}

func AddFridayClip(url string) error {
	db := database.GetDB()

	_, err := db.Exec("INSERT INTO friday_clips (url) VALUES (?)", url)
	return err
}

func DeleteFridayClip(id int64) error {
	db := database.GetDB()

	result, err := db.Exec("DELETE FROM friday_clips WHERE id=?", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Pizza Toppings
func GetPizzaToppings() ([]PizzaTopping, error) {
	db := database.GetDB()

	rows, err := db.Query("SELECT id, name FROM pizza_toppings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var toppings []PizzaTopping
	for rows.Next() {
		var topping PizzaTopping
		err := rows.Scan(&topping.ID, &topping.Name)
		if err != nil {
			return nil, err
		}
		toppings = append(toppings, topping)
	}
	return toppings, nil
}

func AddPizzaTopping(name string) error {
	db := database.GetDB()

	_, err := db.Exec("INSERT INTO pizza_toppings (name) VALUES (?)", name)
	return err
}

func DeletePizzaTopping(id int64) error {
	db := database.GetDB()

	result, err := db.Exec("DELETE FROM pizza_toppings WHERE id=?", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func GetCommandByID(id int64) (*CustomCommand, error) {
	db := database.GetDB()

	var cmd CustomCommand
	err := db.QueryRow("SELECT id, guild_id, prefix, name, description, output, author FROM custom_commands WHERE id = ?", id).
		Scan(&cmd.ID, &cmd.GuildID, &cmd.Prefix, &cmd.Name, &cmd.Description, &cmd.Output, &cmd.Author)
	if err != nil {
		return nil, err
	}
	return &cmd, nil
}
