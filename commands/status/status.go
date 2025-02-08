package status

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
	"wokkibot/utils"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var StatusCommand = discord.SlashCommandCreate{
	Name:        "status",
	Description: "Shows the current status of the bot including version information",
}

func HandleStatus(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		if err := e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, nil); err != nil {
			return err
		}

		self, _ := e.Client().Caches().SelfUser()

		statusEmbed := discord.NewEmbedBuilder().
			SetTitle("Wokkibot Status").
			SetThumbnail(self.EffectiveAvatarURL()).
			AddField("Version", fmt.Sprintf("[%s](https://github.com/Wokkipannu/Wokkibot/commit/%s)", b.Version, b.Version), false).
			AddField("Go", runtime.Version(), false).
			AddField("yt-dlp", getYtdlpVersion(), false).
			AddField("FFmpeg", getFfmpegVersion(), false).
			AddField("Uptime", time.Since(b.StartTime).Round(time.Second).String(), true).
			AddField("Ping", fmt.Sprintf("%dms", b.Client.Gateway().Latency().Milliseconds()), true).
			SetColor(utils.COLOR_GREEN)

		if self.BannerURL() != nil {
			formatOpt := utils.SetCDNOptions(discord.FileFormatPNG, discord.QueryValues{"size": 1024})
			statusEmbed.SetImage(*self.BannerURL(formatOpt))
		}

		_, err := e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetEmbeds(statusEmbed.Build()).
			Build())

		return err
	}
}

func getYtdlpVersion() string {
	cmd := exec.Command("yt-dlp", "--version")
	output, err := cmd.Output()
	ytDlpVersion := "Not found"
	if err == nil {
		ytDlpVersion = strings.TrimSpace(string(output))
	}
	return ytDlpVersion
}

func getFfmpegVersion() string {
	cmd := exec.Command("ffmpeg", "-version")
	output, err := cmd.Output()
	if err != nil {
		return "Not found"
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) == 0 {
		return "Not found"
	}

	parts := strings.Split(lines[0], " ")
	if len(parts) < 3 {
		return "Not found"
	}

	return parts[2]
}
