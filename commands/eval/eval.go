package eval

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	gopiston "github.com/milindmadhukar/go-piston"
)

var discordCodeblockRegex = regexp.MustCompile(`(?s)\x60\x60\x60(?P<language>\w+)\n(?P<code>.+)\x60\x60\x60`)

var EvalCommand = discord.MessageCommandCreate{
	Name: "Eval",
}

func HandleEval(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		msg := e.MessageCommandInteractionData().TargetMessage()

		runtimes, err := b.PistonClient.GetRuntimes()
		if err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Error while getting runtimes").SetFlags(discord.MessageFlagEphemeral).Build())
		}

		matches := discordCodeblockRegex.FindStringSubmatch(msg.Content)
		if len(matches) == 0 {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("No codeblock found").SetFlags(discord.MessageFlagEphemeral).Build())
		}
		rawLanguage := matches[discordCodeblockRegex.SubexpIndex("language")]
		code := matches[discordCodeblockRegex.SubexpIndex("code")]

		var language string
	runtimeloop:
		for _, runtime := range *runtimes {
			if strings.EqualFold(runtime.Language, rawLanguage) {
				language = runtime.Language
				break
			}
			for _, alias := range runtime.Aliases {
				if strings.EqualFold(alias, rawLanguage) {
					language = runtime.Language
					break runtimeloop
				}
			}
		}

		if language == "" {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("Language %s not supported", rawLanguage).SetFlags(discord.MessageFlagEphemeral).Build())
		}

		if err = e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, nil); err != nil {
			return err
		}

		start := time.Now()
		rs, err := b.PistonClient.Execute(language, "", []gopiston.Code{{Content: code}})
		embed := discord.NewEmbedBuilder()
		end := time.Now()
		duration := end.Sub(start)
		if err != nil {
			embed.SetTitle("Eval")
			embed.SetDescriptionf("Error: %s", err.Error())
			embed.AddField("Status", "Error", true)
			embed.AddField("Duration", fmt.Sprintf("%.3f seconds", duration.Seconds()), true)
			embed.AddField("Language", rs.Language, true)
		} else {
			embed.SetTitle("Eval")
			embed.SetDescription(rs.GetOutput())
			embed.AddField("Status", "Success", true)
			embed.AddField("Duration", fmt.Sprintf("%.3f seconds", duration.Seconds()), true)
			embed.AddField("Language", rs.Language, true)
		}

		_, err = e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().SetEmbeds(embed.Build()).AddActionRow(discord.NewLinkButton("View code", msg.JumpURL())).Build())
		return err
	}
}
