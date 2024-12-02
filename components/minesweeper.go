package components

import (
	"wokkibot/commands"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/handler"
)

func HandleMinesweeperFlagAction(b *wokkibot.Wokkibot) handler.ComponentHandler {
	return func(e *handler.ComponentEvent) error {
		commands.HandleMinesweeperFlagAction(b, e)
		return nil
	}
}

func HandleMinesweeperRevealAction(b *wokkibot.Wokkibot) handler.ComponentHandler {
	return func(e *handler.ComponentEvent) error {
		commands.HandleMinesweeperRevealAction(b, e)
		return nil
	}
}

func HandleMinesweeperUpAction(b *wokkibot.Wokkibot) handler.ComponentHandler {
	return func(e *handler.ComponentEvent) error {
		commands.HandleMinesweeperUpAction(b, e)
		return nil
	}
}

func HandleMinesweeperDownAction(b *wokkibot.Wokkibot) handler.ComponentHandler {
	return func(e *handler.ComponentEvent) error {
		commands.HandleMinesweeperDownAction(b, e)
		return nil
	}
}

func HandleMinesweeperLeftAction(b *wokkibot.Wokkibot) handler.ComponentHandler {
	return func(e *handler.ComponentEvent) error {
		commands.HandleMinesweeperLeftAction(b, e)
		return nil
	}
}

func HandleMinesweeperRightAction(b *wokkibot.Wokkibot) handler.ComponentHandler {
	return func(e *handler.ComponentEvent) error {
		commands.HandleMinesweeperRightAction(b, e)
		return nil
	}
}
