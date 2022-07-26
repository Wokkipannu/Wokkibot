package main

import (
	"context"
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
	"github.com/lukasl-dev/waterlink"
)

var (
	session        *discordgo.Session
	token          = config.Config("TOKEN")
	GuildID        = flag.String("guild", config.Config("GUILDID"), "Guild ID for testing slash commands")
	RemoveCommands = flag.Bool("rmcmd", false, "Remove all commands after shutdowning or not")

	userID     = config.Config("USERID")
	passphrase = config.Config("PASSPHRASE")
	host       = config.Config("HOST")

	httpHost, _ = url.Parse(fmt.Sprintf("http://%s", host))
	wsHost, _   = url.Parse(fmt.Sprintf("ws://%s", host))

	connOpts = waterlink.NewConnectOptions().WithUserID(userID).WithPassphrase(passphrase)
	reqOpts  = waterlink.NewRequesterOptions().WithPassphrase(passphrase)

	conn waterlink.Connection
	req  waterlink.Requester

	sessionID string
)

func main() {
	var err error
	session, err = discordgo.New(fmt.Sprintf("Bot %s", token))
	if err != nil {
		log.Fatalln("Error when creating discordgo session:", err)
	}
	log.Println("Bot session created")

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildVoiceStates

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
	wlerr := retry(5, 30*time.Second, func() (err error) {
		werr := initWaterlink()
		return werr
	})
	if wlerr != nil {
		log.Println(wlerr)
		return
	}
	session.UpdateGameStatus(0, "🪗")
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
				conn.Destroy(update.GuildID)
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
	err := conn.UpdateVoice(update.GuildID, sessionID, update.Token, update.Endpoint)
	if err != nil {
		log.Printf("Updating voice server failed on guild %s: %s\n", update.GuildID, err)
	} else {
		log.Printf("Updated voice server of guild %s.\n", update.GuildID)
		commands.Conn = conn
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
	var err error
	conn, err = waterlink.Connect(context.TODO(), *wsHost, connOpts)
	if err != nil {
		return fmt.Errorf("opening connection failed: %v", err)
	}
	commands.Conn = conn
	log.Println("Connection established.")

	req = waterlink.NewRequester(*httpHost, reqOpts)
	commands.Req = req
	log.Println("Requester created.")

	go commands.ListenForEvents()

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
