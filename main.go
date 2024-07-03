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

	customCommands, err := wokkibot.LoadCommands("custom_commands.json")
	if err != nil {
		panic("failed to load custom commands: " + err.Error())
	}

	b := wokkibot.New(*cfg, customCommands)
	defer b.Close()

	r := handler.New()

	r.Command("/ping", commands.HandlePing(b))
	r.Command("/roll", commands.HandleRoll(b))
	r.Command("/flip", commands.HandleFlip(b))
	r.Command("/pizza", commands.HandlePizza(b))
	r.Command("/friday", commands.HandleFriday(b))
	r.Command("/user", commands.HandleUser(b))
	r.Command("/trivia", commands.HandleTrivia(b))
	r.Route("/settings", func(r handler.Router) {
		r.Route("/commands", func(r handler.Router) {
			r.Command("/add", commands.HandleCustomAdd(b))
			r.Command("/remove", commands.HandleCustomRemove(b))
			r.Command("/list", commands.HandleCustomList(b))
		})
		r.Route("/config", func(r handler.Router) {
			r.Command("/system-message", b.AdminMiddleware(commands.HandleAISystemMessageChange(b)))
			r.Command("/model", b.AdminMiddleware(commands.HandleAIModelChange(b)))
			r.Command("/history-count", b.AdminMiddleware(commands.HandleAIHistoryCountChange(b)))
			r.Command("/api_url", b.AdminMiddleware(commands.HandleAIApiUrlChange(b)))
			r.Command("/enabled", b.AdminMiddleware(commands.HandleAIEnableChange(b)))
		})
	})
	r.Command("/joke", commands.HandleJoke(b))
	r.Command("/inspect", commands.HandleInspect(b))
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
