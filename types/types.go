package types

import "github.com/disgoorg/snowflake/v2"

type Command struct {
	Name        string       `json:"name"`
	Prefix      string       `json:"prefix"`
	Description string       `json:"description"`
	Output      string       `json:"output"`
	Author      snowflake.ID `json:"author"`
	GuildID     snowflake.ID `json:"guild_id"`
}
