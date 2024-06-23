package wokkibot

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func (b *Wokkibot) AdminMiddleware(next handler.CommandHandler) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
		// Get the member from the event
		member := event.Member()

		for _, admin := range b.Config.Admins {
			if admin == member.User.ID {
				return next(event)
			}
		}

		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("You do not have permission to use this command").Build())

		return nil
	}
}
