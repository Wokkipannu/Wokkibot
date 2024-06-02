package utils

import (
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

var imageTypes = []string{"image/png", "image/jpeg", "image/gif", "image/webp"}

func QuoteEmbed(msg discord.Message) discord.EmbedBuilder {
	embed := discord.NewEmbedBuilder()
	embed.SetAuthor(msg.Author.Username, "", *msg.Author.AvatarURL())
	embed.SetDescription(msg.Content)
	embed.SetTimestamp(msg.CreatedAt)

	if len(msg.Attachments) > 0 {
		var attachments []string
		for _, attachment := range msg.Attachments {
			a := fmt.Sprintf("[%s](<%s>)", attachment.Filename, attachment.URL)
			attachments = append(attachments, a)
			for _, t := range imageTypes {
				if *attachment.ContentType == t {
					embed.SetImage(attachment.URL)
					break
				}
			}
		}

		embed.AddField("Attachments", strings.Join(attachments, "\n"), true)
	}

	return *embed
}

func FormatDuration(duration lavalink.Duration) string {
	if duration == 0 {
		return "0 minutes 0 seconds"
	}

	minutes := duration.Minutes()
	seconds := duration.SecondsPart()

	minutesText := "minute"
	secondsText := "second"
	if minutes > 1 {
		minutesText += "s"
	}
	if seconds == 1 {
		secondsText += "s"
	}

	return fmt.Sprintf("%d %s %d %s", minutes, minutesText, seconds, secondsText)
}

func FormatPosition(position lavalink.Duration) string {
	if position == 0 {
		return "0:00"
	}
	return fmt.Sprintf("%d:%02d", position.Minutes(), position.SecondsPart())
}
