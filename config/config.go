package config

import (
	"encoding/json"
	"os"

	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"
)

type Config struct {
	Token       string         `json:"token"`
	GuildID     string         `json:"guildid"`
	TriviaToken string         `json:"trivia_token"`
	Admins      []snowflake.ID `json:"admins"`
	Lavalink    LavalinkConfig `json:"lavalink"`
}

type LavalinkConfig struct {
	Enabled bool                   `json:"enabled"`
	Nodes   []disgolink.NodeConfig `json:"nodes"`
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err = json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func SaveConfig(config Config) error {
	file, err := os.OpenFile("config.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Sync()
		_ = file.Close()
	}()
	data, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	return err
}
