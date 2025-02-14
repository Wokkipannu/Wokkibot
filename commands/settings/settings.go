package settings

import (
	"fmt"
	"strings"
	"wokkibot/handlers"
	"wokkibot/types"
	"wokkibot/utils"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/json"
)

var perms = json.NewNullable(discord.Permissions(discord.PermissionAdministrator))

var SettingsCommand = discord.SlashCommandCreate{
	Name:                     "settings",
	Description:              "Used to change server settings",
	DefaultMemberPermissions: &perms,
	Options: []discord.ApplicationCommandOption{
		// Custom command related settings
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
							Choices: []discord.ApplicationCommandOptionChoiceString{
								{
									Name:  "!",
									Value: "!",
								},
								{
									Name:  "?",
									Value: "?",
								},
								{
									Name:  ".",
									Value: ".",
								},
								{
									Name:  ",",
									Value: ",",
								},
								{
									Name:  "-",
									Value: "-",
								},
								{
									Name:  "|",
									Value: "|",
								},
							},
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
		// Guild specific settings, for now only pin channel
		discord.ApplicationCommandOptionSubCommandGroup{
			Name:        "guild",
			Description: "Guild specific settings",
			Options: []discord.ApplicationCommandOptionSubCommand{
				{
					Name:        "pinchannel",
					Description: "Set the pin channel for the guild",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionChannel{
							Name:        "channel",
							Description: "The channel to set as the pin channel",
							Required:    true,
						},
					},
				},
				{
					Name:        "xlinks",
					Description: "Toggles the conversion of x.com links to fixupx.com links",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "state",
							Description: "Enable or disable the conversion of x.com links to fixupx.com links",
							Required:    true,
							Choices: []discord.ApplicationCommandOptionChoiceString{
								{
									Name:  "On",
									Value: "on",
								},
								{
									Name:  "Off",
									Value: "off",
								},
							},
						},
					},
				},
			},
		},
		// Add new lavalink settings group
		discord.ApplicationCommandOptionSubCommandGroup{
			Name:        "lavalink",
			Description: "Manage lavalink settings",
			Options: []discord.ApplicationCommandOptionSubCommand{
				{
					Name:        "toggle",
					Description: "Toggle lavalink on/off",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "state",
							Description: "Turn lavalink on or off",
							Required:    true,
							Choices: []discord.ApplicationCommandOptionChoiceString{
								{
									Name:  "On",
									Value: "on",
								},
								{
									Name:  "Off",
									Value: "off",
								},
							},
						},
					},
				},
			},
		},
	},
}

/**
 * Custom command settings
 */
func HandleCustomAdd(h *handlers.Handler) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		if err := e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, nil); err != nil {
			return err
		}

		data := e.SlashCommandInteractionData()

		prefix := data.String("prefix")
		name := data.String("name")
		description := data.String("description")
		output := data.String("output")

		author := e.User().ID

		err := h.AddOrUpdateCommand(types.Command{
			Name:        name,
			Prefix:      prefix,
			Description: description,
			Output:      output,
			Author:      author,
			GuildID:     *e.GuildID(),
		})
		if err != nil {
			utils.HandleError(e, "Failed to add custom command", err.Error())
			return err
		}

		_, err = e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetContent("Custom command added successfully").
			Build())

		return err
	}
}

func HandleCustomRemove(h *handlers.Handler) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		if err := e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, nil); err != nil {
			return err
		}

		data := e.SlashCommandInteractionData()

		prefix := data.String("prefix")
		name := data.String("name")

		author := e.User().ID

		err := h.RemoveCommand(prefix, name, author)
		if err != nil {
			utils.HandleError(e, "Failed to remove custom command", err.Error())
			return err
		}

		_, err = e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetContent("Custom command removed successfully").
			Build())

		return err
	}
}

func HandleCustomList(h *handlers.Handler) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		var cmds []string
		var descriptions []string
		var authors []string
		for _, cmd := range h.CustomCommands {
			if cmd.GuildID == *e.GuildID() {
				author, _ := e.Client().Rest().GetUser(cmd.Author)
				cmds = append(cmds, fmt.Sprintf("%v%v", cmd.Prefix, cmd.Name))
				descriptions = append(descriptions, cmd.Description)
				authors = append(authors, fmt.Sprintf("%v", author.EffectiveName()))
			}
		}

		embed := discord.NewEmbedBuilder()
		embed.SetTitle("Custom commands")
		embed.SetColor(utils.COLOR_BLURPLE)
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

/**
 * Guild settings
 */
func HandlePinChannelChange(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		if err := e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, nil); err != nil {
			return err
		}

		data := e.SlashCommandInteractionData()

		channel := data.Channel("channel")

		b.Handlers.SetPinChannel(*e.GuildID(), channel.ID)

		_, err := e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetContent("Pin channel updated successfully").
			Build())

		return err
	}
}

func HandleXLinksToggle(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		if err := e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, nil); err != nil {
			return err
		}

		data := e.SlashCommandInteractionData()
		state := data.String("state")

		var message string
		if state == "on" {
			b.Handlers.ToggleGuildXLinks(*e.GuildID(), true)
			message = "X links have been enabled"
		} else {
			b.Handlers.ToggleGuildXLinks(*e.GuildID(), false)
			message = "X links have been disabled"
		}

		_, err := e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetContent(message).
			Build())

		return err
	}
}

/**
 * Lavalink settings
 */
func HandleInitLavalink(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		b.InitLavalink()
		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Lavalink initialized successfully").Build())
	}
}

func HandleLavalinkToggle(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		if err := e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, nil); err != nil {
			return err
		}

		data := e.SlashCommandInteractionData()
		state := data.String("state")

		var message string
		if state == "on" {
			b.InitLavalink()
			message = "Lavalink has been enabled"
		} else {
			if b.Lavalink != nil {
				b.Lavalink.Close()
			}
			b.Config.Lavalink.Enabled = false
			message = "Lavalink has been disabled"
		}

		_, err := e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetContent(message).
			Build())

		return err
	}
}
