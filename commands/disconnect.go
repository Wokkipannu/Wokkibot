package commands

import (
	"wokkibot/utils"

	"github.com/bwmarrin/discordgo"
)

var disconnect = Command{
	Info: &discordgo.ApplicationCommand{
		Name:        "disconnect",
		Description: "Disconnect from channel and destroy the player",
	},
	Run: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		delete(utils.Queue, i.GuildID)
		LeaveVoiceChannel(i.GuildID, i.ChannelID)
		utils.InteractionRespondMessage(s, i, "Queue was deleted and disconnected from voice channel.")
	},
}
