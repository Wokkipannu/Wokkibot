package main

import (
	"log/slog"
	"os"
	"wokkibot/commands"
	"wokkibot/config"
	"wokkibot/database"
	"wokkibot/handlers"
	"wokkibot/queue"
	"wokkibot/web"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
)

var version = "dev"

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// Initialize database
	dbConfig := database.Config{
		DatabaseURL: "file:wokkibot.db",
	}
	if err := database.Initialize(dbConfig); err != nil {
		panic("failed to initialize database: " + err.Error())
	}
	slog.Info("Successfully connected to database")
	defer database.Close()

	// Load custom commands
	customCommands, err := handlers.LoadCommands()
	if err != nil {
		panic("failed to load custom commands: " + err.Error())
	}

	// Load guilds
	guilds, err := handlers.LoadGuilds()
	if err != nil {
		panic("failed to load guilds: " + err.Error())
	}

	// Initialize handlers
	h := handlers.New()

	h.CustomCommands = customCommands

	h.Guilds = guilds

	// Initialize wokkibot
	b := wokkibot.New(*cfg, version, h)
	defer b.Close()

	// Initialize router
	router := handler.New()

	// Initialize queue manager
	queue := queue.NewQueueManager()

	// Register commands
	commands.RegisterCommands(router, b, h, queue)

	// Define intents
	intents := gateway.IntentGuildMessages |
		gateway.IntentDirectMessages |
		gateway.IntentGuildMessageTyping |
		gateway.IntentDirectMessageTyping |
		gateway.IntentMessageContent |
		gateway.IntentGuilds |
		gateway.IntentGuildVoiceStates

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	b.Client, err = disgo.New(b.Config.Token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(intents),
			gateway.WithCompress(true),
			gateway.WithPresenceOpts(
				gateway.WithPlayingActivity("starting up..."),
				gateway.WithOnlineStatus(discord.OnlineStatusDND),
			),
		),
		bot.WithEventListeners(router),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagGuilds, cache.FlagMembers, cache.FlagVoiceStates),
		),
		bot.WithEventListenerFunc(b.OnDiscordEvent),
		bot.WithEventListenerFunc(h.OnMessageCreate),
		bot.WithLogger(logger),
	)

	if err != nil {
		slog.Error("error while building disgo instance", slog.Any("err", err))
		return
	}

	if cfg.Lavalink.Enabled {
		b.InitLavalink()
	}
	if b.Config.GuildID != "" {
		b.SyncGuildCommands(commands.Commands, snowflake.MustParse(b.Config.GuildID))
	} else {
		b.SyncGlobalCommands(commands.Commands)
	}

	// Initialize web server
	webConfig := web.OAuthConfig{
		ClientID:     cfg.Web.ClientID,
		ClientSecret: cfg.Web.ClientSecret,
		RedirectURI:  cfg.Web.RedirectURI,
		AdminUserIDs: cfg.Admins,
	}

	webServer := web.NewServer(webConfig, b, h, version)
	go func() {
		if err := webServer.Start(":3000"); err != nil {
			slog.Error("error starting web server", slog.Any("err", err))
		}
	}()

	b.Start()
}
