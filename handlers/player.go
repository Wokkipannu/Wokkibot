package handlers

import (
	"context"
	"log/slog"
	"wokkibot/queue"

	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

type PlayerHandler struct {
	Queues *queue.QueueManager
}

func NewPlayerHandler(queues *queue.QueueManager) *PlayerHandler {
	return &PlayerHandler{
		Queues: queues,
	}
}

func (h *PlayerHandler) OnPlayerPause(player disgolink.Player, event lavalink.PlayerPauseEvent) {
	slog.Info("player paused", slog.Any("event", event))
}

func (h *PlayerHandler) OnPlayerResume(player disgolink.Player, event lavalink.PlayerResumeEvent) {
	slog.Info("player resumed", slog.Any("event", event))
}

func (h *PlayerHandler) OnTrackStart(player disgolink.Player, event lavalink.TrackStartEvent) {
	slog.Info("track started", slog.Any("event", event))
}

func (h *PlayerHandler) OnTrackEnd(player disgolink.Player, event lavalink.TrackEndEvent) {
	if !event.Reason.MayStartNext() {
		return
	}

	queue := h.Queues.Get(event.GuildID())
	var (
		nextTrack lavalink.Track
		ok        bool
	)

	nextTrack, ok = queue.Next()

	if !ok {
		return
	}
	if err := player.Update(context.TODO(), lavalink.WithTrack(nextTrack)); err != nil {
		slog.Error("Failed to play next track", slog.Any("err", err))
	}
}

func (h *PlayerHandler) OnTrackException(player disgolink.Player, event lavalink.TrackExceptionEvent) {
	slog.Info("track exception", slog.Any("event", event))
}

func (h *PlayerHandler) OnTrackStuck(player disgolink.Player, event lavalink.TrackStuckEvent) {
	slog.Info("track stuck", slog.Any("event", event))
}

func (h *PlayerHandler) OnWebSocketClosed(player disgolink.Player, event lavalink.WebSocketClosedEvent) {
	slog.Info("websocket closed", slog.Any("event", event))
}

func (h *PlayerHandler) OnUnknownEvent(p disgolink.Player, e lavalink.UnknownEvent) {
	slog.Info("unknown event", slog.Any("event", e.Type()), slog.String("data", string(e.Data)))
}
