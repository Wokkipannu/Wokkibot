package commands

import (
	"fmt"
	"log"
	"strings"
	"wokkibot/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/lukasl-dev/waterlink/entity/track"
)

var play = Command{
	Info: &discordgo.ApplicationCommand{
		Name:        "play",
		Description: "Begin playing a track by URL",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "url",
				Description: "Link to the video",
				Required:    true,
			},
		},
	},
	Run: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Member != nil {
			vc := findMembersChannel(i.GuildID, i.Member.User.ID)
			if vc == "" {
				err := utils.InteractionRespondMessage(s, i, "You must be connected to a voice channel")
				if err != nil {
					log.Print(err)
				}
				return
			}

			identifier := i.ApplicationCommandData().Options[0].StringValue()

			track, err := GetTrack(identifier)
			if err != nil {
				utils.InteractionRespondMessage(s, i, fmt.Sprintf("Failed to fetch track: %v", err.Error()))
				return
			}

			if err := s.ChannelVoiceJoinManual(i.GuildID, vc, false, true); err != nil {
				err := utils.InteractionRespondMessage(s, i, "Could not join your voice channel")
				if err != nil {
					log.Print(err)
				}
				return
			}

			if g, ok := utils.Queue[i.GuildID]; ok {
				g.Queue = append(g.Queue, &utils.QueueObj{
					Requester: i.Member,
					Keyword:   identifier,
					TrackInfo: &track.Info,
					TrackID:   track.ID,
				})
			} else {
				newQ := make([]*utils.QueueObj, 1)
				newQ[0] = &utils.QueueObj{
					Requester: i.Member,
					Keyword:   identifier,
					TrackInfo: &track.Info,
					TrackID:   track.ID,
				}
				q := utils.GuildQueue{
					TextChannelID:  i.ChannelID,
					VoiceChannelID: vc,
					Queue:          newQ,
				}
				utils.Queue[i.GuildID] = &q
			}

			joinMemberChannel(vc, i.GuildID, i.Member.User.ID)
			queue := utils.Queue[i.GuildID]
			if len(queue.Queue) > 1 {
				err := utils.InteractionRespondMessage(s, i, "Added to queue")
				if err != nil {
					log.Print(err)
				}
			} else {
				BeginPlay(i.GuildID, i)
			}
		}
	},
}

func GetTrack(identifier string) (*track.Track, error) {
	res, err := Req.LoadTracks(identifier)
	if err != nil {
		return nil, err
	}
	if len(res.Tracks) > 0 {
		track := res.Tracks[0]
		return &track, nil
	} else {
		return nil, fmt.Errorf("search resulted in 0 tracks")
	}
}

func BeginPlay(guildID string, interaction *discordgo.InteractionCreate) {
	q := utils.Queue[guildID]
	if len(q.Queue) == 0 {
		_, _ = Session.ChannelMessageSend(q.TextChannelID, "No more songs in queue")
		delete(utils.Queue, guildID)
		leaveVoiceChannel(guildID, q.TextChannelID)
		return
	}

	if err := Conn.Play(guildID, q.Queue[0].TrackID); err != nil {
		if _, err := Session.ChannelMessageSend(q.TextChannelID, "Could not play track"); err != nil {
			log.Print(err)
		}
		return
	}

	duration := q.Queue[0].TrackInfo.Length
	seconds := (duration / 1000) % 60
	minutes := (duration / (1000 * 60) % 60)
	hours := (duration / (1000 * 60 * 60) % 24)

	color := Session.State.UserColor(Session.State.User.ID, q.TextChannelID)

	embed := &discordgo.MessageEmbed{}
	embed.Color = color
	embed.Title = "Now playing"
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "Title",
		Value:  q.Queue[0].TrackInfo.Title,
		Inline: false,
	})
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "Requester",
		Value:  q.Queue[0].Requester.Nick,
		Inline: true,
	})
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "URL",
		Value:  fmt.Sprintf("[%v](%v)", "Link", q.Queue[0].TrackInfo.URI),
		Inline: true,
	})
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "Duration",
		Value:  fmt.Sprintf("%v:%v:%v", utils.NumberFormat(hours), utils.NumberFormat(minutes), utils.NumberFormat(seconds)),
		Inline: true,
	})

	// If the interaction exists (This function was ran via a command)
	// Send a response to the interaction. If the function was ran via
	// the "TrackEnd", send a normal message with the session
	if interaction != nil {
		if err := utils.InteractionRespondMessageEmbed(Session, interaction, embed); err != nil {
			log.Print(err)
		}
	} else {
		if _, err := Session.ChannelMessageSendEmbed(q.TextChannelID, embed); err != nil {
			log.Print(err)
		}
	}
}

func joinMemberChannel(channelID, guildID, userID string) bool {
	vcID := findMembersChannel(guildID, userID)
	if vcID == "" {
		_, _ = Session.ChannelMessageSend(channelID, "You must be in a voice channel.")
		return false
	}
	if err := Session.ChannelVoiceJoinManual(guildID, vcID, false, true); err != nil {
		_, _ = Session.ChannelMessageSend(channelID, "Could not join your voice channel.")
		return false
	}
	return true
}

func leaveVoiceChannel(guildId, channelId string) bool {
	if err := Session.ChannelVoiceJoinManual(guildId, "", false, true); err != nil {
		_, _ = Session.ChannelMessageSend(channelId, "I was unable to disconnect. Please disconnect me manually.")
		Conn.Destroy(guildId)
		return false
	}
	Conn.Destroy(guildId)
	return true
}

func findMembersChannel(guildID, userID string) string {
	guild, err := Session.State.Guild(guildID)
	if err != nil {
		return ""
	}
	for _, state := range guild.VoiceStates {
		if strings.EqualFold(userID, state.UserID) {
			return state.ChannelID
		}
	}
	return ""
}
