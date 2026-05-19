package eval

import (
	"context"
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

		runtimes, err := b.PistonClient.GetRuntimes(context.Background())
		if err != nil {
			return e.CreateMessage(discord.NewMessageCreate().WithContent("Error while getting runtimes").WithFlags(discord.MessageFlagEphemeral))
		}

		matches := discordCodeblockRegex.FindStringSubmatch(msg.Content)
		if len(matches) == 0 {
			return e.CreateMessage(discord.NewMessageCreate().WithContent("No codeblock found").WithFlags(discord.MessageFlagEphemeral))
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
			return e.CreateMessage(discord.NewMessageCreate().WithContentf("Language %s not supported", rawLanguage).WithFlags(discord.MessageFlagEphemeral))
		}

		if err = e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, nil); err != nil {
			return err
		}

		start := time.Now()
		rs, err := b.PistonClient.Execute(context.Background(), language, "", []gopiston.Code{{Content: code}})
		embed := discord.NewEmbed()
		end := time.Now()
		duration := end.Sub(start)
		if err != nil {
			embed = embed.WithTitle("Eval")
			embed = embed.WithDescriptionf("Error: %s", err.Error())
			embed = embed.AddField("Status", "Error", true)
			embed = embed.AddField("Duration", fmt.Sprintf("%.3f seconds", duration.Seconds()), true)
			embed = embed.AddField("Language", language, true)
		} else {
			embed = embed.WithTitle("Eval")
			embed = embed.WithDescription(rs.GetOutput())
			embed = embed.AddField("Status", "Success", true)
			embed = embed.AddField("Duration", fmt.Sprintf("%.3f seconds", duration.Seconds()), true)
			embed = embed.AddField("Language", rs.Language, true)
		}

		_, err = e.UpdateInteractionResponse(discord.NewMessageUpdate().WithEmbeds(embed).AddActionRow(discord.NewLinkButton("View code", msg.JumpURL())))
		return err
	}
}
