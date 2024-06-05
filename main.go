package main

import (
	"wokkibot/commands"
	"wokkibot/components"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
)

func main() {
	cfg, err := wokkibot.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	b := wokkibot.New(*cfg)
	defer b.Close()

	r := handler.New()
	r.Command("/ping", commands.HandlePing(b))
	r.Command("/roll", commands.HandleRoll(b))
	r.Command("/flip", commands.HandleFlip(b))
	r.Command("/pizza", commands.HandlePizza(b))
	r.Command("/friday", commands.HandleFriday(b))
	r.Command("/user", commands.HandleUser(b))
	r.Command("/trivia", commands.HandleTrivia(b))
	// Context menu commands
	r.Command("/Quote", commands.HandleQuote(b))
	r.Command("/Eval", commands.HandleEval(b))
	// Music commands
	r.Command("/play", commands.HandlePlay(b))
	r.Command("/skip", commands.HandleSkip(b))
	r.Component("/queue/skip", components.HandleQueueSkipAction(b))
	r.Command("/queue", commands.HandleQueue(b))
	r.Command("/disconnect", commands.HandleDisconnect(b))
	r.Command("/seek", commands.HandleSeek(b))
	r.Command("/volume", commands.HandleVolume(b))

	b.SetupBot(r)
	b.InitLavalink()
	if b.Config.GuildID != "" {
		b.SyncGuildCommands(commands.Commands, snowflake.MustParse(b.Config.GuildID))
	} else {
		b.SyncGlobalCommands(commands.Commands)
	}
	b.Start()
}
