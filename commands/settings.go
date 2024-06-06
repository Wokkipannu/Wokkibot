package commands

import (
	"strings"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var settingsCommand = discord.SlashCommandCreate{
	Name:        "settings",
	Description: "Used to change server settings",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommandGroup{
			Name:        "commands",
			Description: "Manage custom commands",
			Options: []discord.ApplicationCommandOptionSubCommand{
				{
					Name:        "add",
					Description: "Add a custom command",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "prefix",
							Description: "The prefix to use for the command",
							Required:    true,
						},
						discord.ApplicationCommandOptionString{
							Name:        "name",
							Description: "The name of the command",
							Required:    true,
						},
						discord.ApplicationCommandOptionString{
							Name:        "description",
							Description: "What does this command do?",
							Required:    true,
						},
						discord.ApplicationCommandOptionString{
							Name:        "output",
							Description: "The output of the command",
							Required:    true,
						},
					},
				},
				{
					Name:        "remove",
					Description: "Remove a custom command",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "prefix",
							Description: "The prefix of the command to remove",
							Required:    true,
						},
						discord.ApplicationCommandOptionString{
							Name:        "name",
							Description: "The name of the command to remove",
							Required:    true,
						},
					},
				},
				{
					Name:        "list",
					Description: "List all custom commands",
				},
			},
		},
	},
}

func HandleCustomAdd(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		data := e.SlashCommandInteractionData()

		prefix := data.String("prefix")
		name := data.String("name")
		description := data.String("description")
		output := data.String("output")

		for i, cmd := range b.CustomCommands {
			if cmd.Prefix == prefix && cmd.Name == name {
				if e.User().ID != cmd.Author {
					return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("You don't have permission to modify %v%v", prefix, name).Build())
				}

				b.CustomCommands[i].Output = output
				b.CustomCommands[i].Prefix = prefix
				b.CustomCommands[i].Name = name
				b.CustomCommands[i].Description = description
				wokkibot.AddOrUpdateCommand("custom_commands.json", b.CustomCommands[i])
				return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("Command **%v%v** modified", prefix, name).Build())
			}
		}

		newCommand := wokkibot.Command{
			Prefix: prefix,
			Name:   name,
			Output: output,
			Author: e.User().ID,
		}

		wokkibot.AddOrUpdateCommand("custom_commands.json", newCommand)

		b.CustomCommands = append(b.CustomCommands, newCommand)

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("Command **%v%v** added", prefix, name).Build())
	}
}

func HandleCustomRemove(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		data := e.SlashCommandInteractionData()

		prefix := data.String("prefix")
		name := data.String("name")

		for _, cmd := range b.CustomCommands {
			if cmd.Prefix == prefix && cmd.Name == name {
				if e.User().ID != cmd.Author {
					return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("You don't have permission to remove %v%v", prefix, name).Build())
				}
			}
		}

		updatedCommands := b.CustomCommands[:0]
		for _, cmd := range b.CustomCommands {
			if !(cmd.Prefix == prefix && cmd.Name == name) {
				b.CustomCommands = append(updatedCommands, cmd)

				wokkibot.RemoveCommand("custom_commands.json", prefix, name)
				return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("Command **%v%v** removed", prefix, name).Build())
			}
		}

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("Command **%v%v** could not be found", prefix, name).Build())
	}
}

func HandleCustomList(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		listCommands := []string{}
		for _, cmd := range b.CustomCommands {
			listCommands = append(listCommands, cmd.Prefix+cmd.Name)
		}
		embed := discord.NewEmbedBuilder()
		embed.SetTitle("Custom commands")
		if len(listCommands) == 0 {
			embed.SetDescription("No custom commands found")
		} else {
			embed.SetDescription(strings.Join(listCommands, "\n"))
		}

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetEmbeds(embed.Build()).Build())
	}
}
