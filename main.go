package main

import (
	"wokkibot/commands"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
)

func main() {
	b := wokkibot.New()

	r := handler.New()
	r.Command("/ping", commands.HandlePing(b))
	r.Command("/roll", commands.HandleRoll(b))
	r.Command("/flip", commands.HandleFlip(b))
	r.Command("/pizza", commands.HandlePizza(b))
	r.Command("/friday", commands.HandleFriday(b))
	r.Command("/user", commands.HandleUser(b))
	// Music commands
	r.Command("/play", commands.HandlePlay(b))
	r.Command("/skip", commands.HandleSkip(b))
	r.Command("/queue", commands.HandleQueue(b))
	r.Command("/disconnect", commands.HandleDisconnect(b))

	b.SetupBot(r)
	b.InitLavalink()
	if wokkibot.Config("GUILDID") != "" {
		b.SyncGuildCommands(commands.Commands, snowflake.MustParse(wokkibot.Config("GUILDID")))
	} else {
		b.SyncGlobalCommands(commands.Commands)
	}
	b.Start()
}
