package commands

import (
	"fmt"
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

		embed := discord.NewEmbedBuilder()
		embed.SetTitle(user.EffectiveName())
		embed.AddField("Nickname", user.Username, false)
		embed.AddField("Joined this server", fmt.Sprintf("%v (%v days ago)", e.Member().JoinedAt.Format("02.01.2006"), DaysSince(e.Member().JoinedAt)), false)
		embed.SetThumbnail(*user.AvatarURL())

		if user.BannerURL() != nil {
			embed.SetImage(*user.BannerURL())
		}

		embed.AddField("Account created", fmt.Sprintf("%v (%v days ago)", user.CreatedAt().Format("02.01.2006"), DaysSince(user.CreatedAt())), false)

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetEmbeds(embed.Build()).Build())
	}
}

func DaysSince(date time.Time) int {
	return int(time.Since(date).Hours() / 24)
}
