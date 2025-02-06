package handlers

import (
	"wokkibot/queue"
	"wokkibot/types"
)

type Handler struct {
	CustomCommands []types.Command
	PlayerHandler  *PlayerHandler
	TriviaManager  *TriviaManager
}

func New() *Handler {
	return &Handler{
		CustomCommands: []types.Command{},
		PlayerHandler:  NewPlayerHandler(queue.NewQueueManager()),
		TriviaManager:  NewTriviaManager(),
	}
}
