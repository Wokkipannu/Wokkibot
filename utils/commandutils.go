package utils

import "github.com/disgoorg/disgo/discord"

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
