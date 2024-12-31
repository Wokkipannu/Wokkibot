package utils

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/disgoorg/disgo/discord"
)

var imageTypes = []string{"image/png", "image/jpeg", "image/gif", "image/webp"}

// Creates a quote embed from a message
func QuoteEmbed(msg discord.Message) discord.EmbedBuilder {
	embed := discord.NewEmbedBuilder()
	embed.SetAuthor(fmt.Sprintf("Quoting %v", msg.Author.EffectiveName()), "", *msg.Author.AvatarURL())
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

// Replaces the domain part of a URL, for example "https://example.com/path" with "https://newdomain.com/path"
func ReplaceDomain(originalURL, newDomain string) (string, error) {
	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		return "", err
	}

	hostParts := strings.Split(newDomain, ":")
	newHost := hostParts[0]
	newPort := ""
	if len(hostParts) > 1 {
		newPort = hostParts[1]
	}

	parsedURL.Host = newHost
	if newPort != "" {
		parsedURL.Host = newHost + ":" + newPort
	}

	return parsedURL.String(), nil
}
