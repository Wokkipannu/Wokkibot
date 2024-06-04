package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
	"time"
	"wokkibot/commands"
	"wokkibot/components"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	dumpGoroutines()
}

func main() {
	http.HandleFunc("/health", healthCheckHandler)
	go http.ListenAndServe(":8080", nil)

	cfg, err := wokkibot.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	b := wokkibot.New(*cfg)
	defer b.Close()

	r := handler.New()
	r.Command("/ping", commands.HandlePing(b))
	r.Command("/roll", commands.HandleRoll(b))
	r.Command("/flip", commands.HandleFlip(b))
	r.Command("/pizza", commands.HandlePizza(b))
	r.Command("/friday", commands.HandleFriday(b))
	r.Command("/user", commands.HandleUser(b))
	r.Command("/trivia", commands.HandleTrivia(b))
	// Context menu commands
	r.Command("/Quote", commands.HandleQuote(b))
	r.Command("/Eval", commands.HandleEval(b))
	// Music commands
	r.Command("/play", commands.HandlePlay(b))
	r.Command("/skip", commands.HandleSkip(b))
	r.Component("/queue/skip", components.HandleQueueSkipAction(b))
	r.Command("/queue", commands.HandleQueue(b))
	r.Command("/disconnect", commands.HandleDisconnect(b))
	r.Command("/seek", commands.HandleSeek(b))
	r.Command("/volume", commands.HandleVolume(b))

	b.SetupBot(r)
	b.InitLavalink()
	if b.Config.GuildID != "" {
		b.SyncGuildCommands(commands.Commands, snowflake.MustParse(b.Config.GuildID))
	} else {
		b.SyncGlobalCommands(commands.Commands)
	}
	b.Start()
}

func dumpGoroutines() {
	timestamp := time.Now().Format("2006-01-02T15-04-05")
	filename := fmt.Sprintf("goroutine_dump_%v.txt", timestamp)
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal("could not create goroutine dump file:", err)
	}
	defer f.Close()
	pprof.Lookup("goroutine").WriteTo(f, 1)
}
