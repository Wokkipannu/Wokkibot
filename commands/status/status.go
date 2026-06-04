package status

import (
	"fmt"
	"runtime"
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

		statusEmbed := createEmbed(b, e)

		_, err := e.UpdateInteractionResponse(discord.NewMessageUpdate().
			WithEmbeds(statusEmbed))

		if err != nil {
			return err
		}

		go func() {
			ping := getPing(b)
			statusEmbed.Fields[6].Value = ping
			_, _ = e.UpdateInteractionResponse(discord.NewMessageUpdate().
				WithEmbeds(statusEmbed))
		}()

		return nil
	}
}

func createEmbed(b *wokkibot.Wokkibot, e *handler.CommandEvent) discord.Embed {
	self, _ := b.Client.Caches.SelfUser()

	currentYtdlpVersion := utils.GetYtdlpVersion()
	latestYtdlpVersion, err := utils.GetLatestYtdlpVersion()
	ytdlpVersion := fmt.Sprintf("%s (Latest: %s)", currentYtdlpVersion, latestYtdlpVersion)
	if err == nil {
		if currentYtdlpVersion == latestYtdlpVersion {
			ytdlpVersion = fmt.Sprintf("%s (Up to date)", currentYtdlpVersion)
		}
	}

	embed := discord.NewEmbed().
		WithTitlef("%s Status", self.Username).
		WithThumbnail(self.EffectiveAvatarURL()).
		AddField("Version", getBotVersion(b), false).
		AddField("Go", runtime.Version(), true).
		AddField("Disgo", utils.GetDisgoVersion(), true).
		AddField("yt-dlp", ytdlpVersion, true).
		AddField("FFmpeg", utils.GetFfmpegVersion(), true).
		AddField("Start time", fmt.Sprintf("<t:%d:R>", b.StartTime.Unix()), true).
		AddField("Ping", getPing(b), true).
		WithColor(utils.COLOR_GREEN)

	if e != nil {
		guild, _ := e.Guild()
		embed = embed.AddField("File Size limit", fmt.Sprintf("%dMB", utils.CalculateMaximumFileSizeForGuild(guild)), true)
	}

	// Self user does not contain BannerURL, so we must fetch it from the client rest
	botUser, err := b.Client.Rest.GetUser(self.ID)

	if err == nil && botUser.BannerURL() != nil {
		formatOpt := utils.SetCDNOptions(discord.FileFormatPNG, discord.QueryValues{"size": 1024})
		embed = embed.WithImage(*botUser.BannerURL(formatOpt))
	}

	return embed
}

func getBotVersion(b *wokkibot.Wokkibot) string {
	if b.Version == "dev" {
		return b.Version
	}

	return fmt.Sprintf("[%s](https://github.com/Wokkipannu/Wokkibot/commit/%s)", b.Version, b.Version)
}

func getPing(b *wokkibot.Wokkibot) string {
	maxRetries := 5
	retryDelay := 2 * time.Second

	for i := range maxRetries {
		ping := b.Client.Gateway.Latency().Milliseconds()
		if ping > 0 {
			return fmt.Sprintf("%dms", ping)
		}
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}
	return "N/A"
}
