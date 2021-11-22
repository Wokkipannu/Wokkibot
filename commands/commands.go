package commands

import (
	"wokkibot/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/lukasl-dev/waterlink"
	"github.com/lukasl-dev/waterlink/entity/event"
	"github.com/lukasl-dev/waterlink/entity/player"
)

type Command struct {
	Info *discordgo.ApplicationCommand
	Run  func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var (
	Commands = []*discordgo.ApplicationCommand{
		play.Info,
		skip.Info,
		volume.Info,
		queue.Info,
		seek.Info,
		// pause.Info,
		// resume.Info,
	}
	Handlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		// Music related commands
		play.Info.Name:   play.Run,
		skip.Info.Name:   skip.Run,
		volume.Info.Name: volume.Run,
		queue.Info.Name:  queue.Run,
		seek.Info.Name:   seek.Run,
		// pause.Info.Name:  pause.Run, // Pause command is not functional
		// resume.Info.Name: resume.Run, // Resume command is not functional
		// Other commands
	}
	Session   *discordgo.Session
	Req       waterlink.Requester
	Conn      waterlink.Connection
	SessionID string
)

func ListenForEvents() {
	for evt := range Conn.Events() {
		switch evt.Type() {
		case event.TrackStuck:
		case event.TrackEnd:
			evt := evt.(player.TrackEnd)
			utils.Queue[evt.GuildID].Queue = utils.Queue[evt.GuildID].Queue[1:]
			BeginPlay(evt.GuildID, nil)
		}
	}
}
