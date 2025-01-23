package commands

import (
	"fmt"
	"strings"
	"wokkibot/utils"
	"wokkibot/wokkibot"

	"wokkibot/database"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/json"
)

var perms = json.NewNullable(discord.Permissions(discord.PermissionAdministrator))

var settingsCommand = discord.SlashCommandCreate{
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
		// Settings for /friday command
		discord.ApplicationCommandOptionSubCommandGroup{
			Name:        "friday",
			Description: "Manage friday celebration clips",
			Options: []discord.ApplicationCommandOptionSubCommand{
				{
					Name:        "add",
					Description: "Add a friday celebration clip",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "url",
							Description: "The URL of the video",
							Required:    true,
						},
					},
				},
				{
					Name:        "remove",
					Description: "Remove a friday celebration clip",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "id",
							Description: "The id of the video",
							Required:    true,
						},
					},
				},
				{
					Name:        "list",
					Description: "List all friday celebration clips",
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
			},
		},
	},
}

func HandleCustomAdd(b *wokkibot.Wokkibot) handler.CommandHandler {
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

		err := b.AddOrUpdateCommand(wokkibot.Command{
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

func HandleCustomRemove(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		if err := e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, nil); err != nil {
			return err
		}

		data := e.SlashCommandInteractionData()

		prefix := data.String("prefix")
		name := data.String("name")

		author := e.User().ID

		err := b.RemoveCommand(prefix, name, author)
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

func HandleAddFridayClip(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		data := e.SlashCommandInteractionData()

		url := data.String("url")

		db := database.GetDB()
		_, err := db.Exec("INSERT INTO friday_clips (url) VALUES (?)", url)
		if err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("Failed to add clip: %v", err).Build())
		}

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Clip added successfully").Build())
	}
}

func HandleRemoveFridayClip(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		data := e.SlashCommandInteractionData()

		url := data.String("id")

		db := database.GetDB()
		_, err := db.Exec("DELETE FROM friday_clips WHERE id = ?", url)
		if err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("Failed to remove clip: %v", err).Build())
		}

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Clip removed successfully").Build())
	}
}

func HandleListFridayClips(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		db := database.GetDB()
		rows, err := db.Query("SELECT id, url FROM friday_clips")
		if err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("Failed to list clips: %v", err).Build())
		}
		defer rows.Close()

		var clips []string
		for rows.Next() {
			var id, url string
			if err := rows.Scan(&id, &url); err != nil {
				return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("Failed to list clips: %v", err).Build())
			}
			clips = append(clips, fmt.Sprintf("%v: <%v>", id, url))
		}
		if err := rows.Err(); err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("Failed to list clips: %v", err).Build())
		}

		if len(clips) == 0 {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("No clips found").Build())
		}

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent(strings.Join(clips, "\n")).Build())
	}
}

func HandlePinChannelChange(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		if err := e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, nil); err != nil {
			return err
		}

		data := e.SlashCommandInteractionData()

		channel := data.Channel("channel")

		db := database.GetDB()
		_, err := db.Exec("UPDATE guilds SET pin_channel = ? WHERE id = ?", channel.ID, *e.GuildID())
		if err != nil {
			utils.HandleError(e, "Failed to update pin channel", err.Error())
			return err
		}

		_, err = e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetContent("Pin channel updated successfully").
			Build())

		return err
	}
}
