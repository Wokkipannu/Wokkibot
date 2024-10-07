package wokkibot

import (
	"encoding/json"
	"os"

	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"
)

type AISettings struct {
	Model        string `json:"model"`
	System       string `json:"system"`
	Enabled      bool   `json:"enabled"`
	ApiUrl       string `json:"api_url"`
	HistoryCount int    `json:"history_count"`
}

type Config struct {
	Token       string                 `json:"token"`
	GuildID     string                 `json:"guildid"`
	Nodes       []disgolink.NodeConfig `json:"nodes"`
	TriviaToken string                 `json:"trivia_token"`
	AISettings  AISettings             `json:"ai_settings"`
	Admins      []snowflake.ID         `json:"admins"`
	PinChannel  snowflake.ID           `json:"pin_channel"`
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
