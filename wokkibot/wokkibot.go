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

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
	gopiston "github.com/milindmadhukar/go-piston"
)

func New(config config.Config, version string) *Wokkibot {
	return &Wokkibot{
		PistonClient: gopiston.CreateDefaultClient(),
		Config:       config,
		Queues: &QueueManager{
			queues: make(map[snowflake.ID]*Queue),
		},
		Trivias: &TriviaManager{
			trivias: make(map[snowflake.ID]*Trivia),
		},
		CustomCommands: []Command{},
		Games:          make(map[snowflake.ID]interface{}),
		StartTime:      time.Now(),
		Version:        version,
	}
}

type Wokkibot struct {
	Client         bot.Client
	Config         config.Config
	PistonClient   *gopiston.Client
	Lavalink       disgolink.Client
	Queues         *QueueManager
	Trivias        *TriviaManager
	CustomCommands []Command
	Games          map[snowflake.ID]interface{}
	StartTime      time.Time
	Version        string
}

func (b *Wokkibot) SetupBot(r handler.Router) {
	var err error
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	b.Client, err = disgo.New(b.Config.Token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentGuildMessages|gateway.IntentDirectMessages|gateway.IntentGuildMessageTyping|gateway.IntentDirectMessageTyping|gateway.IntentMessageContent|gateway.IntentGuilds|gateway.IntentGuildVoiceStates),
			gateway.WithCompress(true),
			gateway.WithPresenceOpts(
				gateway.WithPlayingActivity("starting up..."),
				gateway.WithOnlineStatus(discord.OnlineStatusDND),
			),
		),
		bot.WithEventListeners(r),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagGuilds, cache.FlagMembers, cache.FlagVoiceStates),
		),
		bot.WithEventListenerFunc(b.OnDiscordEvent),
		bot.WithEventListenerFunc(b.onMessageCreate),
		bot.WithLogger(logger),
	)

	if err != nil {
		slog.Error("error while building disgo instance", slog.Any("err", err))
		return
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
		disgolink.WithListenerFunc(b.onPlayerPause),
		disgolink.WithListenerFunc(b.onPlayerResume),
		disgolink.WithListenerFunc(b.onTrackStart),
		disgolink.WithListenerFunc(b.onTrackEnd),
		disgolink.WithListenerFunc(b.onTrackException),
		disgolink.WithListenerFunc(b.onTrackStuck),
		disgolink.WithListenerFunc(b.onWebSocketClosed),
		disgolink.WithListenerFunc(b.onUnknownEvent),
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
			b.Queues.Delete(e.VoiceState.GuildID)
		}
	case *events.VoiceServerUpdate:
		b.Lavalink.OnVoiceServerUpdate(context.TODO(), e.GuildID, e.Token, *e.Endpoint)
	}
}
