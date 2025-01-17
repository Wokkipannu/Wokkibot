package main

import (
	"wokkibot/commands"
	"wokkibot/components"
	"wokkibot/config"
	"wokkibot/database"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
)

var version = "dev"

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	dbConfig := database.Config{
		DatabaseURL: "file:wokkibot.db",
	}
	if err := database.Initialize(dbConfig); err != nil {
		panic("failed to initialize database: " + err.Error())
	}
	defer database.Close()

	b := wokkibot.New(*cfg, []wokkibot.Command{}, version)
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
		r.Route("/friday", func(r handler.Router) {
			r.Command("/add", b.AdminMiddleware(commands.HandleAddFridayClip(b)))
			r.Command("/remove", b.AdminMiddleware(commands.HandleRemoveFridayClip(b)))
			r.Command("/list", b.AdminMiddleware(commands.HandleListFridayClips(b)))
		})
		r.Route("/guild", func(r handler.Router) {
			r.Command("/pinchannel", b.AdminMiddleware(commands.HandlePinChannelChange(b)))
		})
	})
	r.Command("/joke", commands.HandleJoke(b))
	r.Command("/download", commands.HandleDownload(b))
	r.Command("/status", commands.HandleStatus(b))
	// Context menu commands
	r.Command("/Quote", commands.HandleQuote(b))
	r.Command("/Eval", commands.HandleEval(b))
	r.Command("/Pin", commands.HandlePin(b))
	// Music commands
	r.Command("/play", commands.HandlePlay(b))
	r.Command("/skip", commands.HandleSkip(b))
	r.Component("/queue/skip", components.HandleQueueSkipAction(b))
	r.Command("/queue", commands.HandleQueue(b))
	r.Command("/disconnect", commands.HandleDisconnect(b))
	r.Command("/seek", commands.HandleSeek(b))
	r.Command("/volume", commands.HandleVolume(b))
	// Minesweeper
	r.Command("/minesweeper", commands.HandleMinesweeper(b))
	r.Component("/minesweeper/flag", components.HandleMinesweeperFlagAction(b))
	r.Component("/minesweeper/reveal", components.HandleMinesweeperRevealAction(b))
	r.Component("/minesweeper/up", components.HandleMinesweeperUpAction(b))
	r.Component("/minesweeper/down", components.HandleMinesweeperDownAction(b))
	r.Component("/minesweeper/left", components.HandleMinesweeperLeftAction(b))
	r.Component("/minesweeper/right", components.HandleMinesweeperRightAction(b))

	b.SetupBot(r)
	if cfg.LavalinkEnabled {
		b.InitLavalink()
	}
	if b.Config.GuildID != "" {
		b.SyncGuildCommands(commands.Commands, snowflake.MustParse(b.Config.GuildID))
	} else {
		b.SyncGlobalCommands(commands.Commands)
	}
	b.Start()
}
