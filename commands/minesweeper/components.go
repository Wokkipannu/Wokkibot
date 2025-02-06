package minesweeper

import (
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/handler"
)

func HandleMinesweeperFlagActionComponent(b *wokkibot.Wokkibot) handler.ComponentHandler {
	return func(e *handler.ComponentEvent) error {
		HandleMinesweeperFlagAction(b, e)
		return nil
	}
}

func HandleMinesweeperRevealActionComponent(b *wokkibot.Wokkibot) handler.ComponentHandler {
	return func(e *handler.ComponentEvent) error {
		HandleMinesweeperRevealAction(b, e)
		return nil
	}
}

func HandleMinesweeperUpActionComponent(b *wokkibot.Wokkibot) handler.ComponentHandler {
	return func(e *handler.ComponentEvent) error {
		HandleMinesweeperUpAction(b, e)
		return nil
	}
}

func HandleMinesweeperDownActionComponent(b *wokkibot.Wokkibot) handler.ComponentHandler {
	return func(e *handler.ComponentEvent) error {
		HandleMinesweeperDownAction(b, e)
		return nil
	}
}

func HandleMinesweeperLeftActionComponent(b *wokkibot.Wokkibot) handler.ComponentHandler {
	return func(e *handler.ComponentEvent) error {
		HandleMinesweeperLeftAction(b, e)
		return nil
	}
}

func HandleMinesweeperRightActionComponent(b *wokkibot.Wokkibot) handler.ComponentHandler {
	return func(e *handler.ComponentEvent) error {
		HandleMinesweeperRightAction(b, e)
		return nil
	}
}
