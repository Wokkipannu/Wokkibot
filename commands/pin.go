package commands

import (
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var pinCommand = discord.MessageCommandCreate{
	Name: "Pin",
}

func HandlePin(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		msg := e.MessageCommandInteractionData().TargetMessage()

		m := discord.NewMessageCreateBuilder()
		m.SetMessageReference(&discord.MessageReference{
			Type:      discord.MessageReferenceTypeForward,
			MessageID: &msg.ID,
			GuildID:   msg.GuildID,
			ChannelID: &msg.ChannelID,
		})

		_, err := e.Client().Rest().CreateMessage(b.Config.PinChannel, m.Build())
		if err != nil {
			return err
		}

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Message pinned").Build())
	}
}
