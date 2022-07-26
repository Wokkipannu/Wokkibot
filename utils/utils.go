package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"mvdan.cc/xurls/v2"
)

var actualUrl string

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

// Return days since user joined the server
func DaysSince(date time.Time) int {
	return int(time.Since(date).Hours() / 24)
}

//
func GetImageURLFromMessage(m *discordgo.Message) (string, error) {
	attachments := m.Attachments
	if len(attachments) > 0 {
		for _, attachment := range attachments {
			valid, _ := IsValidImage(attachment.URL)
			return valid, nil
		}
	}

	rxStrict := xurls.Strict()
	output := rxStrict.FindAllString(m.Content, -1)
	if len(output) > 0 {
		for _, s := range output {
			valid, _ := IsValidImage(s)
			return valid, nil
		}
	}

	return "", fmt.Errorf("message has no valid images")
}

// Return the given string if it is has a valid suffix
func IsValidImage(search string) (string, error) {
	possibleSuffixes := []string{
		".png",
		".jpg",
		".gif",
		".jpeg",
	}

	for _, s := range possibleSuffixes {
		if strings.HasSuffix(search, s) {
			return search, nil
		}
	}

	// If the link begins with c.tenor.com and ends with .gif, we can just return the search string
	if strings.HasPrefix(search, "https://c.tenor.com/") {
		if strings.HasSuffix(search, ".gif") {
			return search, nil
		}
	}

	// If the string has a prefix to tenor.com/view, we must get the redirect URL from the given URL
	if strings.HasPrefix(search, "https://tenor.com/view/") {
		if !strings.HasSuffix(search, ".gif") {
			search = search + ".gif"
		}

		return GetRedirectURL(search), nil
	}

	return "", fmt.Errorf("could not find any image links")
}

// Get the redirect URL from a request
func GetRedirectURL(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Request error %v", err)
		return ""
	}
	client := new(http.Client)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		actualUrl = req.URL.String()
		return errors.New("Redirect")
	}

	_, err2 := client.Do(req)
	if err2 != nil {
		return actualUrl
	}

	return ""
}
