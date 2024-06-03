package commands

import (
	"fmt"
	"log/slog"
	"strings"
	"time"
	"wokkibot/common"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var userCommand = discord.SlashCommandCreate{
	Name:        "user",
	Description: "Get information about a user",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionUser{
			Name:        "user",
			Description: "The user to get information about",
			Required:    false,
		},
	},
}

func HandleUser(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		data := e.SlashCommandInteractionData()

		var user discord.User

		if u, ok := data.OptUser("user"); ok {
			user = u
		} else {
			user = e.User()
		}

		// For some reason, user does not contain certain attributes, such as BannerURL or AccentColor, so we must fetch the user from the client rest
		fetchedUser, err := b.Client.Rest().GetUser(user.ID)
		if err != nil {
			slog.Info("Error fetching user from client")
		}
		if fetchedUser != nil {
			user = *fetchedUser
		}

		var userFlags []string

		for flag, name := range common.UserFlags {
			if user.PublicFlags&flag != 0 {
				userFlags = append(userFlags, name)
			}
		}

		embed := discord.NewEmbedBuilder()
		embed.SetAuthor(fmt.Sprintf("%v's profile", user.EffectiveName()), "", *user.AvatarURL())
		if user.Bot {
			embed.AddField("Type", "Bot", true)
		} else {
			embed.AddField("Type", "User", true)
		}
		embed.AddField("Global name", user.EffectiveName(), true)
		embed.AddField("Username", user.Username, true)
		if len(userFlags) > 0 {
			embed.AddField("Badges", strings.Join(userFlags, ", "), true)
		}
		embed.AddField("Joined this server", fmt.Sprintf("%v (%v days ago)", e.Member().JoinedAt.Format("02.01.2006"), DaysSince(e.Member().JoinedAt)), false)
		embed.AddField("Account created", fmt.Sprintf("%v (%v days ago)", user.CreatedAt().Format("02.01.2006"), DaysSince(user.CreatedAt())), false)

		embed.SetThumbnail(user.EffectiveAvatarURL())

		if user.AccentColor != nil {
			embed.SetColor(*user.AccentColor)
		}

		if user.BannerURL() != nil {
			formatOpt := SetCDNOptions(discord.FileFormatPNG, discord.QueryValues{"size": 1024})
			embed.SetImage(*user.BannerURL(formatOpt))
		}

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetEmbeds(embed.Build()).Build())
	}
}

func DaysSince(date time.Time) int {
	return int(time.Since(date).Hours() / 24)
}

func SetCDNOptions(format discord.FileFormat, values discord.QueryValues) discord.CDNOpt {
	return func(config *discord.CDNConfig) {
		config.Format = format
		config.Values = values
	}
}
