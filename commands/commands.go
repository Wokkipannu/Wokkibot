package commands

import "github.com/disgoorg/disgo/discord"

var Commands = []discord.ApplicationCommandCreate{
	pingCommand,
	rollCommand,
	flipCommand,
	pizzaCommand,
	fridayCommand,
	userCommand,
	playCommand,
	skipCommand,
	queueCommand,
	disconnectCommand,
}
