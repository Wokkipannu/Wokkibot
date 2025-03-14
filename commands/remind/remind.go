package remind

import (
	"github.com/disgoorg/disgo/discord"
)

var RemindCommand = discord.SlashCommandCreate{
	Name:        "remind",
	Description: "Remind you to do something at a specific time",
	Options: []discord.ApplicationCommandOption{
		SetReminderCommand,
		ListRemindersCommand,
		DeleteReminderCommand,
	},
}
