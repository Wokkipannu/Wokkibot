package commands

import "github.com/disgoorg/disgo/discord"

var Commands = []discord.ApplicationCommandCreate{
	pingCommand,
	rollCommand,
	flipCommand,
	pizzaCommand,
	fridayCommand,
	userCommand,
	quoteCommand,
	evalCommand,
	// Music commands
	playCommand,
	skipCommand,
	queueCommand,
	disconnectCommand,
	seekCommand,
	volumeCommand,
}
