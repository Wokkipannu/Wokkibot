package utils

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lukasl-dev/waterlink/entity/track"
)

type GuildQueue struct {
	TextChannelID  string
	VoiceChannelID string
	Volume         uint
	Queue          []*QueueObj
}

type QueueObj struct {
	Requester *discordgo.Member
	Keyword   string
	TrackInfo *track.Info
	TrackID   string
}

var (
	Queue map[string]*GuildQueue = make(map[string]*GuildQueue)
)
