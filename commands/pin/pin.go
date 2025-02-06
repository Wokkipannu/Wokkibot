package pin

import (
	"wokkibot/database"
	"wokkibot/utils"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
)

var PinCommand = discord.MessageCommandCreate{
	Name: "Pin",
}

func HandlePin(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		if err := e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, nil); err != nil {
			return err
		}

		msg := e.MessageCommandInteractionData().TargetMessage()

		db := database.GetDB()

		var pinChannel string
		err := db.QueryRow("SELECT pin_channel FROM guilds WHERE id = ?", *e.GuildID()).Scan(&pinChannel)
		if err != nil {
			utils.HandleError(e, "Failed to get pin channel", "No pin channel has been set for this guild. Use the `/settings guild pinchannel` command to set one.")
			return err
		}

		m := discord.NewMessageCreateBuilder()
		m.SetMessageReference(&discord.MessageReference{
			Type:      discord.MessageReferenceTypeForward,
			MessageID: &msg.ID,
			GuildID:   msg.GuildID,
			ChannelID: &msg.ChannelID,
		})

		_, err = e.Client().Rest().CreateMessage(snowflake.MustParse(pinChannel), m.Build())
		if err != nil {
			utils.HandleError(e, "Failed to pin message", err.Error())
			return err
		}

		_, err = e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetContent("Message pinned").
			Build())

		return err
	}
}
