package minesweeper

import (
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/handler"
)

func HandleMinesweeperFlagActionComponent(b *wokkibot.Wokkibot) handler.ComponentHandler {
	return func(e *handler.ComponentEvent) error {
		return HandleMinesweeperFlagAction(b, e)
	}
}

func HandleMinesweeperRevealActionComponent(b *wokkibot.Wokkibot) handler.ComponentHandler {
	return func(e *handler.ComponentEvent) error {
		return HandleMinesweeperRevealAction(b, e)
	}
}

func HandleMinesweeperUpActionComponent(b *wokkibot.Wokkibot) handler.ComponentHandler {
	return func(e *handler.ComponentEvent) error {
		return HandleMinesweeperUpAction(b, e)
	}
}

func HandleMinesweeperDownActionComponent(b *wokkibot.Wokkibot) handler.ComponentHandler {
	return func(e *handler.ComponentEvent) error {
		return HandleMinesweeperDownAction(b, e)
	}
}

func HandleMinesweeperLeftActionComponent(b *wokkibot.Wokkibot) handler.ComponentHandler {
	return func(e *handler.ComponentEvent) error {
		return HandleMinesweeperLeftAction(b, e)
	}
}

func HandleMinesweeperRightActionComponent(b *wokkibot.Wokkibot) handler.ComponentHandler {
	return func(e *handler.ComponentEvent) error {
		return HandleMinesweeperRightAction(b, e)
	}
}
