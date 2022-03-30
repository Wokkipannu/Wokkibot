package utils

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// InteractionRespondMessage sends a response to an interaction.
func InteractionRespondMessage(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}

// InteractionRespondMessageEmbed sends a response to an interaction with an embed.
func InteractionRespondMessageEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "",
			Embeds:  []*discordgo.MessageEmbed{embed},
		},
	})
}

// Button returns a discordgo Button
func Button(label string, emoji string, style string, disabled bool, customID string) discordgo.Button {
	var btnStyle discordgo.ButtonStyle
	switch style {
	case "primary":
		btnStyle = discordgo.PrimaryButton
	case "danger":
		btnStyle = discordgo.DangerButton
	case "secondary":
		btnStyle = discordgo.SecondaryButton
	case "success":
		btnStyle = discordgo.SuccessButton
	case "link":
		btnStyle = discordgo.LinkButton
	default:
		btnStyle = discordgo.PrimaryButton
	}

	btn := discordgo.Button{
		Label:    label,
		Style:    btnStyle,
		Disabled: disabled,
		CustomID: customID,
	}
	if emoji != "" {
		btn.Emoji.Name = emoji
	}

	return btn
}

// NumberFormat takes in uint value and adds a leading 0 to the value if it's lower than 10.
func NumberFormat(value uint) string {
	if value < 10 {
		return fmt.Sprintf("0%v", value)
	} else {
		return fmt.Sprintf("%v", value)
	}
}

// IsValidUrl validates given URL and returns a boolean.
func IsValidUrl(uri string) bool {
	_, err := url.ParseRequestURI(uri)
	if err != nil {
		return false
	}
	u, err := url.Parse(uri)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

// TruncateString takes in string and integer and returns truncated string.
func TruncateString(str string, num int) string {
	truncated := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		truncated = str[0:num] + "..."
	}
	return truncated
}

// EscapeString takes in string and replaces its markdown characters to prevent accidental markdown where it should not be.
func EscapeString(str string) string {
	replacer := strings.NewReplacer(
		"[", "",
		"]", "",
		"(", "",
		")", "",
		"*", "",
		"`", "",
		"~", "",
		">", "",
		"||", "",
	)

	str = replacer.Replace(str)

	return str
}

// GetName takes in discordgo Member and returns nickname if the user has one or their username
func GetName(member *discordgo.Member) string {
	if member.Nick != "" {
		return member.Nick
	} else {
		return member.User.Username
	}
}

// Return days since user joined
func DaysSince(member *discordgo.Member) int {
	return int(time.Since(member.JoinedAt).Hours() / 24)
}
