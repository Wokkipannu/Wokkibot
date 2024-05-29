package wokkibot

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"
)

func New() *Wokkibot {
	return &Wokkibot{
		Queues: &QueueManager{
			queues: make(map[snowflake.ID]*Queue),
		},
	}
}

var (
	Token   = Config("TOKEN")
	GuildID = snowflake.MustParse(Config("GUILDID"))

	NodeName      = Config("NODE_NAME")
	NodeAddress   = Config("NODE_ADDRESS")
	NodePassword  = Config("NODE_PASSWORD")
	NodeSecure, _ = strconv.ParseBool(Config("NODE_SECURE"))
)

type Wokkibot struct {
	Client   bot.Client
	Lavalink disgolink.Client
	Queues   *QueueManager
}

func (b *Wokkibot) SetupBot(r handler.Router) {
	var err error
	b.Client, err = disgo.New(Token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentGuildMessages|gateway.IntentDirectMessages|gateway.IntentGuildMessageTyping|gateway.IntentDirectMessageTyping|gateway.IntentMessageContent|gateway.IntentGuilds|gateway.IntentGuildVoiceStates),
			gateway.WithPresenceOpts(
				gateway.WithPlayingActivity("🤡"),
				gateway.WithOnlineStatus(discord.OnlineStatusDND),
			),
		),
		bot.WithEventListeners(r),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagVoiceStates),
		),
		bot.WithEventListenerFunc(b.onVoiceStateUpdate),
		bot.WithEventListenerFunc(b.onVoiceServerUpdate),
	)

	if err != nil {
		slog.Error("error while building disgo instance", slog.Any("err", err))
		return
	}

}

func (b *Wokkibot) SyncGuildCommands(commands []discord.ApplicationCommandCreate) {
	if _, err := b.Client.Rest().SetGuildCommands(b.Client.ApplicationID(), GuildID, commands); err != nil {
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	node, err := b.Lavalink.AddNode(ctx, disgolink.NodeConfig{
		Name:     NodeName,
		Address:  NodeAddress,
		Password: NodePassword,
		Secure:   NodeSecure,
	})

	if err != nil {
		slog.Error("error while adding lavalink node", slog.Any("err", err))
	}

	version, err := node.Version(ctx)
	if err != nil {
		slog.Error("error while getting lavalink version", slog.Any("err", err))
	}

	slog.Info("Lavalink connection established", slog.String("node_version", version), slog.String("node_session_id", node.SessionID()))
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

func (b *Wokkibot) onVoiceStateUpdate(event *events.GuildVoiceStateUpdate) {
	if event.VoiceState.UserID != b.Client.ApplicationID() {
		return
	}
	b.Lavalink.OnVoiceStateUpdate(context.TODO(), event.VoiceState.GuildID, event.VoiceState.ChannelID, event.VoiceState.SessionID)
	if event.VoiceState.ChannelID == nil {
		b.Queues.Delete(event.VoiceState.GuildID)
	}
}

func (b *Wokkibot) onVoiceServerUpdate(event *events.VoiceServerUpdate) {
	b.Lavalink.OnVoiceServerUpdate(context.TODO(), event.GuildID, event.Token, *event.Endpoint)
}
