package handlers

import (
	"github.com/disgoorg/snowflake/v2"
)

type BlackjackGame struct {
	IsActive bool
}

func (g *BlackjackGame) SetStatus(status bool) {
	g.IsActive = status
}

type BlackjackManager struct {
	games map[snowflake.ID]*BlackjackGame
}

func NewBlackjackManager() *BlackjackManager {
	return &BlackjackManager{
		games: make(map[snowflake.ID]*BlackjackGame),
	}
}

func (m *BlackjackManager) Get(guildID snowflake.ID) *BlackjackGame {
	game, ok := m.games[guildID]
	if !ok {
		game = &BlackjackGame{
			IsActive: false,
		}
		m.games[guildID] = game
	}
	return game
}
