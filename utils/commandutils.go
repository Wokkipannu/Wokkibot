package utils

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

func QuoteEmbed(msg discord.Message) discord.EmbedBuilder {
	embed := discord.NewEmbedBuilder()
	embed.SetAuthor(msg.Author.Username, "", *msg.Author.AvatarURL())
	embed.SetDescription(msg.Content)
	embed.SetTimestamp(msg.CreatedAt)

	if len(msg.Attachments) > 0 {
		for _, attachment := range msg.Attachments {
			embed.SetImage(attachment.URL)
		}
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
