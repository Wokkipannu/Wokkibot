package wokkibot

import (
	"context"
	"log/slog"

	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

func (b *Wokkibot) onPlayerPause(player disgolink.Player, event lavalink.PlayerPauseEvent) {
	slog.Info("player paused", slog.Any("event", event))
}

func (b *Wokkibot) onPlayerResume(player disgolink.Player, event lavalink.PlayerResumeEvent) {
	slog.Info("player resumed", slog.Any("event", event))
}

func (b *Wokkibot) onTrackStart(player disgolink.Player, event lavalink.TrackStartEvent) {
	slog.Info("track started", slog.Any("event", event))
}

func (b *Wokkibot) onTrackEnd(player disgolink.Player, event lavalink.TrackEndEvent) {
	if !event.Reason.MayStartNext() {
		return
	}

	queue := b.Queues.Get(event.GuildID())
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

func (b *Wokkibot) onTrackException(player disgolink.Player, event lavalink.TrackExceptionEvent) {
	slog.Info("track exception", slog.Any("event", event))
}

func (b *Wokkibot) onTrackStuck(player disgolink.Player, event lavalink.TrackStuckEvent) {
	slog.Info("track stuck", slog.Any("event", event))
}

func (b *Wokkibot) onWebSocketClosed(player disgolink.Player, event lavalink.WebSocketClosedEvent) {
	slog.Info("websocket closed", slog.Any("event", event))
}

func (b *Wokkibot) onUnknownEvent(p disgolink.Player, e lavalink.UnknownEvent) {
	slog.Info("unknown event", slog.Any("event", e.Type()), slog.String("data", string(e.Data)))
}
