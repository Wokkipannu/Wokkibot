package commands

import (
	"fmt"
	"log/slog"
	"time"
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

		// For some reason, user does not contain certain attributes, such as BannerURL or AccentColor, so we must fetch the user from the client
		fetchedUser, err := b.Client.Rest().GetUser(user.ID)
		if err != nil {
			slog.Info("Error fetching user from client")
		}

		embed := discord.NewEmbedBuilder()
		if fetchedUser.Bot {
			embed.SetTitlef("%v 🤖", fetchedUser.EffectiveName())
		} else {
			embed.SetTitle(fetchedUser.EffectiveName())
		}
		embed.AddField("Nickname", *fetchedUser.GlobalName, false)
		embed.AddField("Joined this server", fmt.Sprintf("%v (%v days ago)", e.Member().JoinedAt.Format("02.01.2006"), DaysSince(e.Member().JoinedAt)), false)
		embed.SetThumbnail(fetchedUser.EffectiveAvatarURL())

		if fetchedUser.AccentColor != nil {
			embed.SetColor(*fetchedUser.AccentColor)
		}

		if fetchedUser.BannerURL() != nil {
			formatOpt := SetCDNOptions(discord.FileFormatPNG, discord.QueryValues{"size": 4092})
			embed.SetImage(*fetchedUser.BannerURL(formatOpt))
		}

		embed.AddField("Account created", fmt.Sprintf("%v (%v days ago)", fetchedUser.CreatedAt().Format("02.01.2006"), DaysSince(fetchedUser.CreatedAt())), false)

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
