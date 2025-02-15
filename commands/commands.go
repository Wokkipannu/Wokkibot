package commands

import (
	"wokkibot/commands/download"
	"wokkibot/commands/eval"
	"wokkibot/commands/flip"
	"wokkibot/commands/friday"
	"wokkibot/commands/joke"
	"wokkibot/commands/minesweeper"
	"wokkibot/commands/music"
	"wokkibot/commands/pin"
	"wokkibot/commands/ping"
	"wokkibot/commands/pizza"
	"wokkibot/commands/quote"
	"wokkibot/commands/remind"
	"wokkibot/commands/roll"
	"wokkibot/commands/settings"
	"wokkibot/commands/status"
	"wokkibot/commands/trivia"
	"wokkibot/commands/user"
	"wokkibot/handlers"
	"wokkibot/middleware"
	"wokkibot/queue"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var Commands = []discord.ApplicationCommandCreate{
	ping.PingCommand,
	roll.RollCommand,
	flip.FlipCommand,
	pizza.PizzaCommand,
	friday.FridayCommand,
	user.UserCommand,
	quote.QuoteCommand,
	eval.EvalCommand,
	trivia.TriviaCommand,
	settings.SettingsCommand,
	joke.JokeCommand,
	pin.PinCommand,
	download.DownloadCommand,
	minesweeper.MinesweeperCommand,
	status.StatusCommand,
	remind.RemindCommand,
	// Music commands
	music.PlayCommand,
	music.SkipCommand,
	music.QueueCommand,
	music.DisconnectCommand,
	music.SeekCommand,
	music.VolumeCommand,
}

func RegisterCommands(r *handler.Mux, b *wokkibot.Wokkibot, h *handlers.Handler, q *queue.QueueManager) {
	r.Command("/ping", ping.HandlePing(b))
	r.Command("/roll", roll.HandleRoll(b))
	r.Command("/flip", flip.HandleFlip(b))
	r.Command("/pizza", pizza.HandlePizza(b))
	r.Command("/friday", friday.HandleFriday(b))
	r.Command("/user", user.HandleUser(b))
	r.Command("/trivia", trivia.HandleTrivia(b))
	r.Route("/settings", func(r handler.Router) {
		r.Route("/commands", func(r handler.Router) {
			r.Command("/add", settings.HandleCustomAdd(h))
			r.Command("/remove", settings.HandleCustomRemove(h))
			r.Command("/list", settings.HandleCustomList(h))
		})
		// r.Route("/friday", func(r handler.Router) {
		// 	r.Command("/add", middleware.AdminMiddleware(settings.HandleAddFridayClip(b)))
		// 	r.Command("/remove", middleware.AdminMiddleware(settings.HandleRemoveFridayClip(b)))
		// 	r.Command("/list", middleware.AdminMiddleware(settings.HandleListFridayClips(b)))
		// })
		r.Route("/guild", func(r handler.Router) {
			r.Command("/pinchannel", middleware.AdminMiddleware(settings.HandlePinChannelChange(b)))
			r.Command("/xlinks", middleware.AdminMiddleware(settings.HandleXLinksToggle(b)))
		})
		r.Route("/lavalink", func(r handler.Router) {
			r.Command("/toggle", middleware.AdminMiddleware(settings.HandleLavalinkToggle(b)))
		})
	})
	r.Command("/joke", joke.HandleJoke(b))
	r.Command("/download", download.HandleDownload(b))
	r.Command("/status", status.HandleStatus(b))
	r.Command("/remind", remind.HandleRemind(b))
	// Context menu commands
	r.Command("/Quote", quote.HandleQuote(b))
	r.Command("/Eval", eval.HandleEval(b))
	r.Command("/Pin", pin.HandlePin(b))
	// Music commands
	r.Command("/play", music.HandlePlay(b, q))
	r.Command("/skip", music.HandleSkip(b, q))
	r.Component("/queue/skip", music.HandleQueueSkipActionComponent(b, q))
	r.Command("/queue", music.HandleQueue(b, q))
	r.Command("/disconnect", music.HandleDisconnect(b))
	r.Command("/seek", music.HandleSeek(b))
	r.Command("/volume", music.HandleVolume(b))
	// Minesweeper
	r.Command("/minesweeper", minesweeper.HandleMinesweeper(b))
	r.Component("/minesweeper/flag", minesweeper.HandleMinesweeperFlagActionComponent(b))
	r.Component("/minesweeper/reveal", minesweeper.HandleMinesweeperRevealActionComponent(b))
	r.Component("/minesweeper/up", minesweeper.HandleMinesweeperUpActionComponent(b))
	r.Component("/minesweeper/down", minesweeper.HandleMinesweeperDownActionComponent(b))
	r.Component("/minesweeper/left", minesweeper.HandleMinesweeperLeftActionComponent(b))
	r.Component("/minesweeper/right", minesweeper.HandleMinesweeperRightActionComponent(b))
}
