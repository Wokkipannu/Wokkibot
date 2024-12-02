package commands

import (
	"fmt"
	"runtime"
	"time"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var statusCommand = discord.SlashCommandCreate{
	Name:        "status",
	Description: "Shows the current status of the bot including version information",
}

func HandleStatus(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		statusEmbed := discord.NewEmbedBuilder().
			SetTitle("Wokkibot Status").
			AddField("Running Version", b.Version, true).
			AddField("Go Version", runtime.Version(), true).
			AddField("Uptime", time.Since(b.StartTime).Round(time.Second).String(), true).
			AddField("Ping", fmt.Sprintf("%dms", b.Client.Gateway().Latency().Milliseconds()), true)

		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetEmbeds(statusEmbed.Build()).
			Build())
	}
}
