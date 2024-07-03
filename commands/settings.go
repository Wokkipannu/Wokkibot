package commands

import (
	"fmt"
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
		discord.ApplicationCommandOptionSubCommandGroup{
			Name:        "config",
			Description: "Change bots configuration",
			Options: []discord.ApplicationCommandOptionSubCommand{
				{
					Name:        "system-message",
					Description: "Change the system message",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "system_message",
							Description: "The system message",
							Required:    true,
						},
					},
				},
				{
					Name:        "model",
					Description: "Change the model",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "model",
							Description: "The model to use",
							Required:    true,
						},
					},
				},
				{
					Name:        "history-count",
					Description: "Change the history count",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionInt{
							Name:        "history_count",
							Description: "The history count",
							Required:    true,
						},
					},
				},
				{
					Name:        "api_url",
					Description: "Change the API URL",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "api_url",
							Description: "The API URL",
							Required:    true,
						},
					},
				},
				{
					Name:        "enabled",
					Description: "Change the enabled status",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionBool{
							Name:        "enabled",
							Description: "The enabled status",
							Required:    true,
						},
					},
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
			if cmd.Prefix == prefix && cmd.Name == name && cmd.GuildID == *e.GuildID() {
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
			Prefix:      prefix,
			Name:        name,
			Description: description,
			Output:      output,
			Author:      e.User().ID,
			GuildID:     *e.GuildID(),
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

		var commandToRemove *wokkibot.Command
		for _, cmd := range b.CustomCommands {
			if cmd.Prefix == prefix && cmd.Name == name && cmd.GuildID == *e.GuildID() {
				if e.User().ID != cmd.Author {
					return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("You don't have permission to remove %v%v", prefix, name).Build())
				}
				commandToRemove = &cmd
				break
			}
		}

		if commandToRemove != nil {
			updatedCommands := make([]wokkibot.Command, 0, len(b.CustomCommands))
			for _, cmd := range b.CustomCommands {
				if !(cmd.Prefix == prefix && cmd.Name == name && cmd.GuildID == *e.GuildID()) {
					updatedCommands = append(updatedCommands, cmd)
				}
			}
			b.CustomCommands = updatedCommands
			wokkibot.RemoveCommand("custom_commands.json", prefix, name)
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("Command **%v%v** removed", prefix, name).Build())
		}

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("Command **%v%v** could not be found", prefix, name).Build())
	}
}

func HandleCustomList(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		var cmds []string
		var descriptions []string
		var authors []string
		for _, cmd := range b.CustomCommands {
			if cmd.GuildID == *e.GuildID() {
				author, _ := b.Client.Rest().GetUser(cmd.Author)
				cmds = append(cmds, fmt.Sprintf("%v%v", cmd.Prefix, cmd.Name))
				descriptions = append(descriptions, cmd.Description)
				authors = append(authors, fmt.Sprintf("%v", author.EffectiveName()))
			}
		}

		embed := discord.NewEmbedBuilder()
		embed.SetTitle("Custom commands")
		if len(cmds) == 0 {
			embed.SetDescription("No custom commands found")
		} else {
			embed.AddField("Command", strings.Join(cmds, "\n"), true)
			embed.AddField("Description", strings.Join(descriptions, "\n"), true)
			embed.AddField("Author", strings.Join(authors, "\n"), true)
		}

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetEmbeds(embed.Build()).Build())
	}
}

func HandleAISystemMessageChange(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		data := e.SlashCommandInteractionData()

		systemMessage := data.String("system_message")

		b.Config.AISettings.System = systemMessage
		wokkibot.SaveConfig(b.Config)

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("System message updated").Build())
	}
}

func HandleAIModelChange(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		data := e.SlashCommandInteractionData()

		model := data.String("model")

		b.Config.AISettings.Model = model
		wokkibot.SaveConfig(b.Config)

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("Model set to %v", model).Build())
	}
}

func HandleAIHistoryCountChange(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		data := e.SlashCommandInteractionData()

		historyCount := data.Int("history_count")

		b.Config.AISettings.HistoryCount = historyCount
		wokkibot.SaveConfig(b.Config)

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("History count set to %v", historyCount).Build())
	}
}

func HandleAIApiUrlChange(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		data := e.SlashCommandInteractionData()

		apiUrl := data.String("api_url")

		b.Config.AISettings.ApiUrl = apiUrl
		wokkibot.SaveConfig(b.Config)

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("API URL set to %v", apiUrl).Build())
	}
}

func HandleAIEnableChange(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		data := e.SlashCommandInteractionData()

		enabled := data.Bool("enabled")

		b.Config.AISettings.Enabled = enabled
		wokkibot.SaveConfig(b.Config)

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("AI enabled set to %v", enabled).Build())
	}
}
