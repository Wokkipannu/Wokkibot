package main

import (
	"wokkibot/commands"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/handler"
)

func main() {
	b := wokkibot.New()

	r := handler.New()
	r.Command("/ping", commands.HandlePing(b))
	r.Command("/roll", commands.HandleRoll(b))
	r.Command("/flip", commands.HandleFlip(b))
	r.Command("/pizza", commands.HandlePizza(b))
	r.Command("/friday", commands.HandleFriday(b))
	// Music commands
	r.Command("/play", commands.HandlePlay(b))
	r.Command("/skip", commands.HandleSkip(b))

	b.SetupBot(r)
	b.InitLavalink()
	b.SyncGuildCommands(commands.Commands)
	b.Start()
}
