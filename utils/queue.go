package utils

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lukasl-dev/waterlink/v2/track"
)

type GuildQueue struct {
	TextChannelID  string
	VoiceChannelID string
	Volume         uint
	Queue          []*QueueObj
}

type QueueObj struct {
	Requester   *discordgo.Member
	Keyword     string
	Track       *track.Track
	TrackInfo   *track.Info
	TrackID     string
	Interaction *discordgo.Interaction
	Message     *discordgo.Message
	Embed       *discordgo.MessageEmbed
}

var (
	Queue map[string]*GuildQueue = make(map[string]*GuildQueue)
)
