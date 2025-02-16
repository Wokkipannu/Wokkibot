package types

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type Command struct {
	Name        string       `json:"name"`
	Prefix      string       `json:"prefix"`
	Description string       `json:"description"`
	Output      string       `json:"output"`
	Author      snowflake.ID `json:"author"`
	GuildID     snowflake.ID `json:"guild_id"`
}

type Guild struct {
	ID            snowflake.ID `json:"id"`
	PinChannel    snowflake.ID `json:"pin_channel"`
	TriviaToken   string       `json:"trivia_token"`
	ConvertXLinks bool         `json:"convert_x_links"`
}

type Reminder struct {
	ID        int          `json:"id"`
	UserID    snowflake.ID `json:"user_id"`
	ChannelID snowflake.ID `json:"channel_id"`
	GuildID   snowflake.ID `json:"guild_id"`
	Message   string       `json:"message"`
	RemindAt  time.Time    `json:"remind_at"`
}
