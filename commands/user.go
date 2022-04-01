package commands

import (
	"fmt"
	"log"
	"wokkibot/utils"

	"github.com/bwmarrin/discordgo"
)

var user = Command{
	Info: &discordgo.ApplicationCommand{
		Name:        "user",
		Description: "Get information about a user",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionMentionable,
				Name:        "user",
				Description: "User to get information about",
				Required:    true,
			},
		},
	},
	Run: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		user := i.ApplicationCommandData().Options[0].UserValue(s)
		if user == nil {
			if err := utils.InteractionRespondMessage(s, i, "User not found"); err != nil {
				log.Print(err)
			}
			return
		}

		member, err := s.GuildMember(i.GuildID, user.ID)
		if err != nil {
			if err := utils.InteractionRespondMessage(s, i, "User not found"); err != nil {
				log.Print(err)
			}
			return
		}

		embed := &discordgo.MessageEmbed{
			Title: fmt.Sprintf("%s#%s", user.Username, user.Discriminator),
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Nickname",
					Value:  utils.GetName(member),
					Inline: true,
				},
				{
					Name:   "Joined this server",
					Value:  fmt.Sprintf("%v (%v days ago)", member.JoinedAt.Format("02.01.2006"), utils.DaysSince(member.JoinedAt)),
					Inline: true,
				},
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: user.AvatarURL(""),
			},
		}
		if user.BannerURL("") != "" {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Banner",
				Value:  fmt.Sprintf("[Click here](%s)", user.BannerURL("4096")),
				Inline: true,
			})
		}

		snowflake, err := discordgo.SnowflakeTimestamp(user.ID)
		if err == nil {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Account created",
				Value:  fmt.Sprintf("%v (%v days ago)", snowflake.Format("02.01.2006"), utils.DaysSince(snowflake)),
				Inline: false,
			})
		}

		embed.Color = user.AccentColor

		if err := utils.InteractionRespondMessageEmbed(s, i, embed); err != nil {
			log.Print(err)
		}
	},
}
