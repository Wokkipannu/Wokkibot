package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"wokkibot/commands"
	"wokkibot/config"
	"wokkibot/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/gompus/snowflake"
	"github.com/lukasl-dev/waterlink/v2"
)

var (
	session        *discordgo.Session
	token          = flag.String("token", config.Config("TOKEN"), "Discord bot account token")
	GuildID        = flag.String("guild", config.Config("GUILDID"), "Guild ID for testing slash commands")
	RemoveCommands = flag.Bool("rmcmd", false, "Remove all commands after shutdowning or not")

	userID     = flag.String("userid", config.Config("USERID"), "ID of the discord bot account")
	passphrase = flag.String("passphrase", config.Config("PASSPHRASE"), "Lavalink passphrase")
	host       = flag.String("host", config.Config("HOST"), "Lavalink host")

	httpHost, _ = url.Parse(fmt.Sprintf("http://%s", *host))
	wsHost, _   = url.Parse(fmt.Sprintf("ws://%s", *host))

	WaterlinkClient     *waterlink.Client
	WaterlinkConnection waterlink.Connection

	sessionID string
)

func main() {
	var err error
	session, err = discordgo.New(fmt.Sprintf("Bot %s", *token))
	if err != nil {
		log.Fatalln("Error when creating discordgo session:", err)
	}
	log.Println("Bot session created")

	session.Identify.Intents = discordgo.IntentsMessageContent | discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildVoiceStates

	session.AddHandler(handleReady)
	session.AddHandler(handleInteractionCreate)
	session.AddHandler(handleVoiceUpdate)
	session.AddHandler(handleVoiceStateUpdate)
	session.AddHandler(handleMessageCreate)

	if err := session.Open(); err != nil {
		log.Fatalln("Error when opening discordgo session:", err)
	}
	log.Println("Bot session opened")

	cmds := registerCommands()
	// wlerr := retry(5, 30*time.Second, func() (err error) {
	// 	werr := initWaterlink()
	// 	return werr
	// })
	// if wlerr != nil {
	// 	log.Println(wlerr)
	// 	return
	// }
	session.UpdateGameStatus(0, "😎")
	commands.Session = session

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	log.Println("Gracefully shutting down")

	if *RemoveCommands {
		for _, cmd := range cmds {
			err := session.ApplicationCommandDelete(session.State.User.ID, *GuildID, cmd.ID)
			if err != nil {
				log.Fatalf("Cannot delete %q command: %v", cmd.Name, err)
			}
		}
	}

	session.Close()
}

func handleReady(_ *discordgo.Session, ready *discordgo.Ready) {
	sessionID = ready.SessionID
	commands.SessionID = ready.SessionID
}

func handleVoiceStateUpdate(s *discordgo.Session, update *discordgo.VoiceStateUpdate) {
	if s.State.User.ID == update.UserID {
		if update.ChannelID == "" {
			if q, ok := utils.Queue[update.GuildID]; ok {
				WaterlinkConnection.Guild(snowflake.MustParse(update.GuildID)).Destroy()
				commands.LeaveVoiceChannel(update.GuildID, q.VoiceChannelID)
				delete(utils.Queue, update.GuildID)
			}
		} else {
			if q, ok := utils.Queue[update.GuildID]; ok {
				q.VoiceChannelID = update.ChannelID
			}
		}
	}
}

func handleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if h, ok := commands.Handlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
			log.Printf("Command %v was ran", i.ApplicationCommandData().Name)
		}
	case discordgo.InteractionMessageComponent:
		if h, ok := commands.Handlers[i.MessageComponentData().CustomID]; ok {
			h(s, i)
			log.Printf("Component %v was ran", i.MessageComponentData().CustomID)
		}
	}
}

func handleVoiceUpdate(_ *discordgo.Session, update *discordgo.VoiceServerUpdate) {
	g := WaterlinkConnection.Guild(snowflake.MustParse(update.GuildID))
	err := g.UpdateVoice(update.GuildID, update.Token, update.Endpoint)
	if err != nil {
		log.Printf("Updating voice server failed on guild %s: %s\n", update.GuildID, err)
	} else {
		log.Printf("Updated voice server of guild %s.\n", update.GuildID)
	}
}

func registerCommands() []*discordgo.ApplicationCommand {
	cmds, err := session.ApplicationCommandBulkOverwrite(session.State.User.ID, *GuildID, commands.Commands)
	if err != nil {
		log.Fatalf("Cannot register commands: %v", err)
	}
	return cmds
}

func initWaterlink() (wlerr error) {
	creds := waterlink.Credentials{
		Authorization: *passphrase,
		UserID:        snowflake.MustParse(session.State.User.ID),
	}

	// Create waterlink client
	WaterlinkClient, err := waterlink.NewClient(httpHost.String(), creds)
	if err != nil {
		return fmt.Errorf("creating client failed: %v", err)
	}
	commands.WaterlinkClient = WaterlinkClient
	log.Println("Waterlink client created")

	// Create waterlink connection
	WaterlinkConnection, err := waterlink.Open(wsHost.String(), creds)
	if err != nil {
		return fmt.Errorf("opening connection failed: %v", err)
	}

	commands.WaterlinkConnection = WaterlinkConnection
	log.Println("Waterlnk connection established")

	// conn, err = waterlink.Connect(context.TODO(), *wsHost, connOpts)
	// if err != nil {
	// 	return fmt.Errorf("opening connection failed: %v", err)
	// }
	// commands.Conn = conn
	// log.Println("Connection established.")

	// req = waterlink.NewRequester(*httpHost, reqOpts)
	// commands.Req = req
	// log.Println("Requester created.")

	// go commands.ListenForEvents()

	return nil
}

func retry(attempts int, sleep time.Duration, f func() error) (err error) {
	for i := 0; ; i++ {
		err = f()
		if err == nil {
			return
		}

		if i >= (attempts - 1) {
			break
		}

		time.Sleep(sleep)

		log.Println("retrying after error:", err)
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}

// Message Create listener has only one purpose. To check if the message includes a link to a message.
// If it does, the bot will attempt to create an embed with the content of the message.
// This only works if the message is in the same channel the link was posted in.
func handleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	prefix := "https://discord.com/channels/"

	// log.Printf("Message from %v: %v", m.Author.Username, m.Content)

	// if m.MessageReference != nil {
	// 	log.Printf("has messagereference")
	// }

	// if m.ReferencedMessage != nil {
	// 	log.Printf("has referencedmessage")
	// }

	// // Log out the message components
	// for _, c := range m.Components {
	// 	log.Printf("Component: %v", c)
	// }

	// // Log out the embeds in the message
	// for _, e := range m.Embeds {
	// 	log.Printf("Embed: %v", e)
	// }

	// if m.MentionChannels != nil {
	// 	log.Printf("has mentionchannels")
	// }

	// // Log out all mentions in the message
	// for _, u := range m.Mentions {
	// 	log.Printf("Mention: %v", u)
	// }

	if strings.HasPrefix(m.Content, prefix) {
		slashes := strings.Split(m.Content, "/")
		messageId := slashes[len(slashes)-1]
		msg, err := s.ChannelMessage(m.ChannelID, messageId)
		if err != nil {
			return
		}

		embed := &discordgo.MessageEmbed{}
		embed.Author = &discordgo.MessageEmbedAuthor{
			Name:    "Quote from " + msg.Author.Username,
			IconURL: msg.Author.AvatarURL(""),
		}
		embed.Description = msg.Content
		embed.Timestamp = msg.Timestamp.Format(time.RFC3339)
		embed.Color = msg.Author.AccentColor

		img, imgErr := utils.GetImageURLFromMessage(msg)
		if imgErr == nil {
			embed.Image = &discordgo.MessageEmbedImage{
				URL: img,
			}
		}

		messageData := &discordgo.MessageSend{
			Content: "",
			Embeds:  []*discordgo.MessageEmbed{embed},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label: "Go to message",
							Style: discordgo.LinkButton,
							URL:   "https://discord.com/channels/" + m.GuildID + "/" + msg.ChannelID + "/" + msg.ID,
							Emoji: discordgo.ComponentEmoji{
								Name: "🔗",
							},
						},
					},
				},
			},
		}

		_, err = s.ChannelMessageSendComplex(m.ChannelID, messageData)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
