package handlers

import (
	"database/sql"
	"fmt"
	"log/slog"
	"wokkibot/database"
	"wokkibot/types"

	"github.com/disgoorg/snowflake/v2"
)

func LoadGuilds() (map[snowflake.ID]types.Guild, error) {
	db := database.GetDB()
	rows, err := db.Query("SELECT id, trivia_token, convert_x_links FROM guilds")
	if err != nil {
		return nil, fmt.Errorf("failed to query guilds: %v", err)
	}
	defer rows.Close()

	guilds := make(map[snowflake.ID]types.Guild)
	for rows.Next() {
		var guild types.Guild
		var idStr string
		var triviaTokenStr sql.NullString
		if err := rows.Scan(&idStr, &triviaTokenStr, &guild.ConvertXLinks); err != nil {
			return nil, fmt.Errorf("failed to scan guild row: %v", err)
		}

		guild.ID = snowflake.MustParse(idStr)
		if triviaTokenStr.Valid {
			guild.TriviaToken = triviaTokenStr.String
		}

		guilds[guild.ID] = guild
	}

	return guilds, nil
}

func (h *Handler) EnsureGuildExists(guildID snowflake.ID) {
	if _, exists := h.Guilds[guildID]; !exists {
		h.Guilds[guildID] = types.Guild{
			ID:            guildID,
			ConvertXLinks: true,
		}

		db := database.GetDB()
		_, err := db.Exec("INSERT INTO guilds (id, convert_x_links) VALUES (?, ?) ON CONFLICT(id) DO NOTHING",
			guildID.String(), true)
		if err != nil {
			slog.Error("Failed to create default guild settings", "guild_id", guildID, "error", err)
		}
	}
}

func (h *Handler) ToggleGuildXLinks(guildID snowflake.ID, enabled bool) {
	if guildID != 0 {
		h.EnsureGuildExists(guildID)
	}

	db := database.GetDB()

	guild, exists := h.Guilds[guildID]
	if !exists {
		return
	}

	guild.ConvertXLinks = enabled
	h.Guilds[guildID] = guild

	_, err := db.Exec("UPDATE guilds SET convert_x_links = $1 WHERE id = $2", enabled, guildID)
	if err != nil {
		slog.Error("error toggling guild x links", slog.Any("err", err))
	}
}
