package commands

import (
	"math/rand"
	"time"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var fridayCommand = discord.SlashCommandCreate{
	Name:        "friday",
	Description: "Post a friday celebration video",
}

func HandleFriday(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		var videos [5]string
		videos[0] = "https://cdn.discordapp.com/attachments/754470348145295360/908680252157480980/fonkymonkyfriday.mp4"
		videos[1] = "https://cdn.discordapp.com/attachments/754470348145295360/908671890111991848/fonky_monky_2.mp4"
		videos[2] = "https://cdn.discordapp.com/attachments/754470348145295360/975746878777987082/perjantai.mp4"
		videos[3] = "https://cdn.discordapp.com/attachments/754470348145295360/975746878371151882/nyt_on_perjantai.mp4"
		videos[4] = "https://cdn.discordapp.com/attachments/754470348145295360/975746876110409809/Perjantai_1.mp4"

		var video int

		r := rand.NewSource(time.Now().UnixNano())
		video = rand.New(r).Intn(len(videos))

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent(videos[video]).Build())
	}
}
