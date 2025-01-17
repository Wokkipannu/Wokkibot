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
		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Custom commands are getting a rewrite").Build())

		data := e.SlashCommandInteractionData()

		prefix := data.String("prefix")
		name := data.String("name")
		description := data.String("description")
		output := data.String("output")

		db := database.GetDB()

		var existingAuthor string
		err := db.QueryRow("SELECT author FROM custom_commands WHERE guild_id = ? AND prefix = ? AND name = ?",
			*e.GuildID(), prefix, name).Scan(&existingAuthor)

		if err == nil {
			if existingAuthor != e.User().ID.String() {
				return e.CreateMessage(discord.NewMessageCreateBuilder().
					SetContentf("You don't have permission to modify %v%v", prefix, name).Build())
			}

			_, err = db.Exec(`
                UPDATE custom_commands 
                SET description = ?, output = ? 
                WHERE guild_id = ? AND prefix = ? AND name = ?`,
				description, output, *e.GuildID(), prefix, name)

			if err != nil {
				return e.CreateMessage(discord.NewMessageCreateBuilder().
					SetContentf("Failed to update command: %v", err).Build())
			}

			return e.CreateMessage(discord.NewMessageCreateBuilder().
				SetContentf("Command **%v%v** updated successfully!", prefix, name).Build())
		}

		_, err = db.Exec(`
            INSERT INTO custom_commands (guild_id, prefix, name, description, output, author) 
            VALUES (?, ?, ?, ?, ?, ?)`,
			*e.GuildID(), prefix, name, description, output, e.User().ID)

		if err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().
				SetContentf("Failed to create command: %v", err).Build())
		}

		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContentf("Command **%v%v** created successfully!", prefix, name).Build())

		// for i, cmd := range b.CustomCommands {
		// 	if cmd.Prefix == prefix && cmd.Name == name && cmd.GuildID == *e.GuildID() {
		// 		if e.User().ID != cmd.Author {
		// 			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("You don't have permission to modify %v%v", prefix, name).Build())
		// 		}

		// 		b.CustomCommands[i].Output = output
		// 		b.CustomCommands[i].Prefix = prefix
		// 		b.CustomCommands[i].Name = name
		// 		b.CustomCommands[i].Description = description
		// 		wokkibot.AddOrUpdateCommand("custom_commands.json", b.CustomCommands[i])
		// 		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("Command **%v%v** modified", prefix, name).Build())
		// 	}
		// }

		// newCommand := wokkibot.Command{
		// 	Prefix:      prefix,
		// 	Name:        name,
		// 	Description: description,
		// 	Output:      output,
		// 	Author:      e.User().ID,
		// 	GuildID:     *e.GuildID(),
		// }

		// wokkibot.AddOrUpdateCommand("custom_commands.json", newCommand)

		// b.CustomCommands = append(b.CustomCommands, newCommand)

		// return e.CreateMessage(discord.NewMessageCreateBuilder().SetContentf("Command **%v%v** added", prefix, name).Build())
	}
}

func HandleCustomRemove(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Custom commands are getting a rewrite").Build())

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
		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Custom commands are getting a rewrite").Build())

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
