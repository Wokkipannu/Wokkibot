package handlers

import (
	"wokkibot/queue"
	"wokkibot/types"

	"github.com/disgoorg/snowflake/v2"
)

type Handler struct {
	CustomCommands []types.Command
	Guilds         map[snowflake.ID]types.Guild
	PlayerHandler  *PlayerHandler
	TriviaManager  *TriviaManager
}

func New() *Handler {
	return &Handler{
		CustomCommands: []types.Command{},
		Guilds:         make(map[snowflake.ID]types.Guild),
		PlayerHandler:  NewPlayerHandler(queue.NewQueueManager()),
		TriviaManager:  NewTriviaManager(),
	}
}
