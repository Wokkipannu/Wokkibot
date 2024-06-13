package components

import (
	"wokkibot/commands"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/handler"
)

func HandleQueueSkipAction(b *wokkibot.Wokkibot) handler.ComponentHandler {
	return func(e *handler.ComponentEvent) error {
		commands.HandleQueueSkipAction(b, e)
		return nil
	}
}
