package components

import (
	"log/slog"
	"wokkibot/commands"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/handler"
)

func HandleQueueSkipAction(b *wokkibot.Wokkibot) handler.ComponentHandler {
	return func(e *handler.ComponentEvent) error {
		slog.Info("queue skip action")
		commands.HandleQueueSkipAction(b, e)
		return nil
	}
}
