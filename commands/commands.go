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
	triviaCommand,
	settingsCommand,
	jokeCommand,
	pinCommand,
	downloadCommand,
	minesweeperCommand,
	statusCommand,
	// Music commands
	playCommand,
	skipCommand,
	queueCommand,
	disconnectCommand,
	seekCommand,
	volumeCommand,
}
