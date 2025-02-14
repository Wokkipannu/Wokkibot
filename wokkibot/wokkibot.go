package wokkibot

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"wokkibot/config"
	"wokkibot/handlers"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
	gopiston "github.com/milindmadhukar/go-piston"
)

type Wokkibot struct {
	Client       bot.Client
	Config       config.Config
	PistonClient *gopiston.Client
	Lavalink     disgolink.Client
	Trivias      *handlers.TriviaManager
	Games        map[snowflake.ID]interface{}
	StartTime    time.Time
	Version      string
	Handlers     *handlers.Handler
}

func New(config config.Config, version string, handlers *handlers.Handler) *Wokkibot {
	return &Wokkibot{
		PistonClient: gopiston.CreateDefaultClient(),
		Config:       config,
		Handlers:     handlers,
		Trivias:      handlers.TriviaManager,
		Games:        make(map[snowflake.ID]interface{}),
		StartTime:    time.Now(),
		Version:      version,
	}
}

func (b *Wokkibot) SyncGuildCommands(commands []discord.ApplicationCommandCreate, guildID snowflake.ID) {
	if _, err := b.Client.Rest().SetGuildCommands(b.Client.ApplicationID(), guildID, commands); err != nil {
		slog.Error("error while registering guild commands", slog.Any("err", err))
	}
}

func (b *Wokkibot) SyncGlobalCommands(commands []discord.ApplicationCommandCreate) {
	if _, err := b.Client.Rest().SetGlobalCommands(b.Client.ApplicationID(), commands); err != nil {
		slog.Error("error while registering global commands", slog.Any("err", err))
	}
}

func (b *Wokkibot) InitLavalink() {
	b.Lavalink = disgolink.New(b.Client.ApplicationID(),
		disgolink.WithListenerFunc(b.Handlers.PlayerHandler.OnPlayerPause),
		disgolink.WithListenerFunc(b.Handlers.PlayerHandler.OnPlayerResume),
		disgolink.WithListenerFunc(b.Handlers.PlayerHandler.OnTrackStart),
		disgolink.WithListenerFunc(b.Handlers.PlayerHandler.OnTrackEnd),
		disgolink.WithListenerFunc(b.Handlers.PlayerHandler.OnTrackException),
		disgolink.WithListenerFunc(b.Handlers.PlayerHandler.OnTrackStuck),
		disgolink.WithListenerFunc(b.Handlers.PlayerHandler.OnWebSocketClosed),
		disgolink.WithListenerFunc(b.Handlers.PlayerHandler.OnUnknownEvent),
	)

	var wg sync.WaitGroup
	for i := range b.Config.Lavalink.Nodes {
		wg.Add(1)
		cfg := b.Config.Lavalink.Nodes[i]
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			node, err := b.Lavalink.AddNode(ctx, cfg)
			if err != nil {
				slog.Error("error while adding lavalink node", slog.Any("err", err))
				b.Config.Lavalink.Enabled = false
				return
			}

			if err = node.Update(context.Background(), lavalink.SessionUpdate{
				Resuming: json.Ptr(true),
				Timeout:  json.Ptr(100),
			}); err != nil {
				slog.Error("error while updating lavalink node", slog.Any("err", err))
			}

			version, err := node.Version(ctx)
			if err != nil {
				slog.Error("error while getting lavalink version", slog.Any("err", err))
			}

			slog.Info("Lavalink connection established", slog.String("node_version", version), slog.String("node_session_id", node.SessionID()))
			b.Config.Lavalink.Enabled = true
		}()
	}
	wg.Wait()

	b.Config.Lavalink.Enabled = true
}

func (b *Wokkibot) Start() {
	if err := b.Client.OpenGateway(context.TODO()); err != nil {
		slog.Error("error while opening gateway", slog.Any("err", err))
	}

	slog.Info("Wokkibot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}

func (b *Wokkibot) Close() {
	b.Lavalink.ForNodes(func(node disgolink.Node) {
		for i, cfgNode := range b.Config.Lavalink.Nodes {
			if node.Config().Name == cfgNode.Name {
				b.Config.Lavalink.Nodes[i].SessionID = node.SessionID()
			}
		}
	})

	b.Lavalink.Close()
	b.Client.Close(context.Background())
}

func (b *Wokkibot) OnDiscordEvent(event bot.Event) {
	switch e := event.(type) {
	case *events.Ready:
		if err := b.Client.SetPresence(context.TODO(), gateway.WithListeningActivity("Bobr kurwa 🦫"), gateway.WithOnlineStatus(discord.OnlineStatusOnline)); err != nil {
			slog.Error("error while setting presence", slog.Any("err", err))
		}
	case *events.GuildVoiceStateUpdate:
		if e.VoiceState.UserID != b.Client.ApplicationID() {
			return
		}
		b.Lavalink.OnVoiceStateUpdate(context.TODO(), e.VoiceState.GuildID, e.VoiceState.ChannelID, e.VoiceState.SessionID)
		if e.VoiceState.ChannelID == nil {
			b.Handlers.PlayerHandler.Queues.Delete(e.VoiceState.GuildID)
		}
	case *events.VoiceServerUpdate:
		b.Lavalink.OnVoiceServerUpdate(context.TODO(), e.GuildID, e.Token, *e.Endpoint)
	}
}
