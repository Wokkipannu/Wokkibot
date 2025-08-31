package utils

import (
	"fmt"
	"net/url"
	"strings"
	"net/http"
	"os/exec"
	"runtime/debug"
	"encoding/json"
	"os"
	"wokkibot/database"

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

// Sets the CDN options for a URL
func SetCDNOptions(format discord.FileFormat, values discord.QueryValues) discord.CDNOpt {
	return func(config *discord.CDNConfig) {
		config.Format = format
		config.Values = values
	}
}

func UpdateStatistics(statsitic string) {
	db := database.GetDB()
	db.Exec("UPDATE statistics SET " + statsitic + " = " + statsitic + " + 1")
}

func GetYtdlpVersion() string {
    cmd := exec.Command("yt-dlp", "--version")
    output, err := cmd.Output()
    ytdlpVersion := "Not found"
    if err == nil {
        ytdlpVersion = strings.TrimSpace(string(output))
    }
    return ytdlpVersion
}

func GetLatestYtdlpVersion() (string, error) {
    resp, err := http.Get("https://api.github.com/repos/yt-dlp/yt-dlp/releases/latest")
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var release struct {
        TagName string `json:"tag_name"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
        return "", err
    }
    return strings.TrimPrefix(release.TagName, ""), nil
}

func UpdateYtdlpBinary() error {
    curlCmd := exec.Command("curl", "-L", "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp", "-o", "/usr/local/bin/yt-dlp")
    if out, err := curlCmd.CombinedOutput(); err != nil {
        return fmt.Errorf("curl update failed: %w, output: %s", err, string(out))
    }

    if err := os.Chmod("/usr/local/bin/yt-dlp", 0755); err != nil {
        return fmt.Errorf("chmod failed: %w", err)
    }

    return nil
}

func GetFfmpegVersion() string {
	cmd := exec.Command("ffmpeg", "-version")
	output, err := cmd.Output()
	if err != nil {
		return "Not found"
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) == 0 {
		return "Not found"
	}

	parts := strings.Split(lines[0], " ")
	if len(parts) < 3 {
		return "Not found"
	}

	return parts[2]
}

func GetDisgoVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, dep := range info.Deps {
			if dep.Path == "github.com/disgoorg/disgo" {
				return dep.Version
			}
		}
	}
	return "Unknown"
}