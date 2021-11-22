package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"wokkibot/commands"
	"wokkibot/config"

	"github.com/bwmarrin/discordgo"
	"github.com/lukasl-dev/waterlink"
)

var (
	session        *discordgo.Session
	token          = config.Config("TOKEN")
	GuildID        = flag.String("guild", config.Config("GUILDID"), "Guild ID for testing slash commands")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")

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

	if err := session.Open(); err != nil {
		log.Fatalln("Error when opening discordgo session:", err)
	}
	log.Println("Bot session opened")

	cmds := registerCommands()
	initWaterlink()
	go commands.ListenForEvents() // Listen for waterlink events
	// session.UpdateGameStatus(0, "")
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

func handleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if h, ok := commands.Handlers[i.ApplicationCommandData().Name]; ok {
		h(s, i)

		log.Printf("Command %v was ran", i.ApplicationCommandData().Name)
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
	// for _, v := range commands.Commands {
	// 	_, err := session.ApplicationCommandCreate(session.State.User.ID, *GuildID, v)
	// 	if err != nil {
	// 		log.Printf("Cannot create '%v' command: '%v'", v.Name, err)
	// 	}
	// }
}

func initWaterlink() {
	var err error
	conn, err = waterlink.Connect(context.TODO(), *wsHost, connOpts)
	if err != nil {
		log.Fatalln("Opening connection failed:", err)
	}
	commands.Conn = conn
	log.Println("Connection established.")

	req = waterlink.NewRequester(*httpHost, reqOpts)
	commands.Req = req
	log.Println("Requester created.")
}
