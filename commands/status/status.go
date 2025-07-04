package status

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
	"wokkibot/database"
	"wokkibot/types"
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

		statusEmbed := createEmbed(b, e, nil)

		_, err := e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetEmbeds(statusEmbed.Build()).
			Build())

		if err != nil {
			return err
		}

		go func() {
			ping := getPing(b)
			statusEmbed.Fields[6].Value = ping
			_, _ = e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
				SetEmbeds(statusEmbed.Build()).
				AddActionRow(discord.NewPrimaryButton("Statistics", "/status/statistics").WithEmoji(discord.ComponentEmoji{Name: "📊"})).
				Build())
		}()

		return nil
	}
}

func createEmbed(b *wokkibot.Wokkibot, e *handler.CommandEvent, c *handler.ComponentEvent) *discord.EmbedBuilder {
	self, _ := b.Client.Caches().SelfUser()

	currentYtdlpVersion := getYtdlpVersion()
	latestYtdlpVersion, err := getLatestYtdlpVersion()
	ytdlpVersion := fmt.Sprintf("%s (Latest: %s)", currentYtdlpVersion, latestYtdlpVersion)
	if err == nil {
		if currentYtdlpVersion == latestYtdlpVersion {
			ytdlpVersion = fmt.Sprintf("%s (Up to date)", currentYtdlpVersion)
		}
	}

	embed := discord.NewEmbedBuilder().
		SetTitlef("%s Status", self.Username).
		SetThumbnail(self.EffectiveAvatarURL()).
		AddField("Version", fmt.Sprintf("[%s](https://github.com/Wokkipannu/Wokkibot/commit/%s)", b.Version, b.Version), false).
		AddField("Go", runtime.Version(), true).
		AddField("Disgo", getDisgoVersion(), true).
		AddField("yt-dlp", ytdlpVersion, true).
		AddField("FFmpeg", getFfmpegVersion(), true).
		AddField("Uptime", fmt.Sprintf("<t:%d:R>", b.StartTime.Unix()), true).
		AddField("Ping", getPing(b), true).
		SetColor(utils.COLOR_GREEN)

	if c != nil {
		guild, _ := c.Guild()
		embed.AddField("File Size limit", fmt.Sprintf("%dMB", utils.CalculateMaximumFileSizeForGuild(guild)), true)
	}

	if e != nil {
		guild, _ := e.Guild()
		embed.AddField("File Size limit", fmt.Sprintf("%dMB", utils.CalculateMaximumFileSizeForGuild(guild)), true)
	}

	// Self user does not contain BannerURL, so we must fetch it from the client rest
	botUser, err := b.Client.Rest().GetUser(self.ID)

	if err == nil && botUser.BannerURL() != nil {
		formatOpt := utils.SetCDNOptions(discord.FileFormatPNG, discord.QueryValues{"size": 1024})
		embed.SetImage(*botUser.BannerURL(formatOpt))
	}

	return embed
}

func getYtdlpVersion() string {
	cmd := exec.Command("yt-dlp", "--version")
	output, err := cmd.Output()
	ytdlpVersion := "Not found"
	if err == nil {
		ytdlpVersion = strings.TrimSpace(string(output))
	}
	return ytdlpVersion
}

func getLatestYtdlpVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/yt-dlp/yt-dlp/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}
	return strings.TrimPrefix(release.TagName, ""), nil
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

func getDisgoVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, dep := range info.Deps {
			if dep.Path == "github.com/disgoorg/disgo" {
				return dep.Version
			}
		}
	}
	return "Unknown"
}

func getPing(b *wokkibot.Wokkibot) string {
	maxRetries := 5
	retryDelay := 2 * time.Second

	for i := range maxRetries {
		ping := b.Client.Gateway().Latency().Milliseconds()
		if ping > 0 {
			return fmt.Sprintf("%dms", ping)
		}
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}
	return "N/A"
}

func HandleStatusStatistics(b *wokkibot.Wokkibot) handler.ComponentHandler {
	return func(e *handler.ComponentEvent) error {
		self, _ := b.Client.Caches().SelfUser()

		db := database.GetDB()
		var statistics types.Statistics
		err := db.QueryRow("SELECT video_downloads, names_given, songs_played, pizzas_generated, coins_flipped, dice_rolled, trivia_games_played, trivia_games_won, trivia_games_lost FROM statistics").Scan(&statistics.VideoDownloads, &statistics.NamesGiven, &statistics.SongsPlayed, &statistics.PizzasGenerated, &statistics.CoinsFlipped, &statistics.DiceRolled, &statistics.TriviaGamesPlayed, &statistics.TriviaGamesWon, &statistics.TriviaGamesLost)
		if err != nil {
			return err
		}

		embed := discord.NewEmbedBuilder().
			SetTitlef("%s Statistics", self.Username).
			SetThumbnail(self.EffectiveAvatarURL()).
			AddField("Video Downloads", fmt.Sprintf("%d", statistics.VideoDownloads), true).
			AddField("Names Given", fmt.Sprintf("%d", statistics.NamesGiven), true).
			AddField("Songs Played", fmt.Sprintf("%d", statistics.SongsPlayed), true).
			AddField("Pizzas Generated", fmt.Sprintf("%d", statistics.PizzasGenerated), true).
			AddField("Coins Flipped", fmt.Sprintf("%d", statistics.CoinsFlipped), true).
			AddField("Dice Rolled", fmt.Sprintf("%d", statistics.DiceRolled), true).
			AddField("Trivia Games Played", fmt.Sprintf("%d", statistics.TriviaGamesPlayed), true).
			AddField("Trivia Games Won", fmt.Sprintf("%d", statistics.TriviaGamesWon), true).
			AddField("Trivia Games Lost", fmt.Sprintf("%d", statistics.TriviaGamesLost), true).
			SetColor(utils.COLOR_GREEN)

		err = e.Respond(discord.InteractionResponseTypeUpdateMessage, discord.NewMessageUpdateBuilder().
			SetEmbeds(embed.Build()).
			AddActionRow(discord.NewPrimaryButton("Status", "/status/status").WithEmoji(discord.ComponentEmoji{Name: "📺"}).WithDisabled(true)).
			Build())

		if err != nil {
			return err
		}

		go func() {
			time.Sleep(5 * time.Second)
			_, _ = e.Client().Rest().UpdateMessage(e.Channel().ID(), e.Message.ID, discord.NewMessageUpdateBuilder().
				SetEmbeds(embed.Build()).
				AddActionRow(discord.NewPrimaryButton("Status", "/status/status").WithEmoji(discord.ComponentEmoji{Name: "📺"}).WithDisabled(false)).
				Build())
		}()

		return nil
	}
}

func HandleStatusStatus(b *wokkibot.Wokkibot) handler.ComponentHandler {
	return func(e *handler.ComponentEvent) error {
		statusEmbed := createEmbed(b, nil, e)

		err := e.Respond(discord.InteractionResponseTypeUpdateMessage, discord.NewMessageUpdateBuilder().
			SetEmbeds(statusEmbed.Build()).
			AddActionRow(discord.NewPrimaryButton("Statistics", "/status/statistics").WithEmoji(discord.ComponentEmoji{Name: "📊"}).WithDisabled(true)).
			Build())

		if err != nil {
			return err
		}

		go func() {
			time.Sleep(5 * time.Second)
			_, _ = e.Client().Rest().UpdateMessage(e.Channel().ID(), e.Message.ID, discord.NewMessageUpdateBuilder().
				SetEmbeds(statusEmbed.Build()).
				AddActionRow(discord.NewPrimaryButton("Statistics", "/status/statistics").WithEmoji(discord.ComponentEmoji{Name: "📊"}).WithDisabled(false)).
				Build())
		}()

		return nil
	}
}
