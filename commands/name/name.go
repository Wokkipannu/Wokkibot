package name

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"wokkibot/utils"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var NameCommand = discord.SlashCommandCreate{
	Name:        "name",
	Description: "Generates a random two-part name from the names list",
}

func Init() {}

func HandleName(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		if err := e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, nil); err != nil {
			return err
		}

		cmdDir, err := os.Getwd()
		if err != nil {
			utils.HandleError(e, "Error while getting current directory", err.Error())
			return err
		}

		namesFile := filepath.Join(cmdDir, "names.txt")

		file, err := os.Open(namesFile)
		if err != nil {
			utils.HandleError(e, "Error while opening names file", err.Error())
			return err
		}
		defer file.Close()

		var names []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			name := strings.TrimSpace(scanner.Text())
			if name != "" {
				names = append(names, name)
			}
		}

		if err := scanner.Err(); err != nil {
			utils.HandleError(e, "Error while scanning names file", err.Error())
			return err
		}

		if len(names) < 2 {
			utils.HandleError(e, "Not enough names in the file", "")
			return err
		}

		firstIndex := rand.Intn(len(names))
		secondIndex := rand.Intn(len(names))

		for secondIndex == firstIndex {
			secondIndex = rand.Intn(len(names))
		}

		randomName := fmt.Sprintf("You are **%s%s**", names[firstIndex], names[secondIndex])

		_, err = e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetContent(randomName).
			Build())

		return err
	}
}
