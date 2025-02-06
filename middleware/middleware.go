package middleware

import (
	"wokkibot/config"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func AdminMiddleware(next handler.CommandHandler) handler.CommandHandler {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	return func(event *handler.CommandEvent) error {
		member := event.Member()

		for _, admin := range cfg.Admins {
			if admin == member.User.ID {
				return next(event)
			}
		}

		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("You do not have permission to use this command").Build())

		return nil
	}
}
