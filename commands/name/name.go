package name

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
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

		namesLen := big.NewInt(int64(len(names)))
		firstBig, err := rand.Int(rand.Reader, namesLen)
		if err != nil {
			utils.HandleError(e, "Error generating random number", err.Error())
			return err
		}
		firstIndex := int(firstBig.Int64())

		var secondIndex int
		for {
			secondBig, err := rand.Int(rand.Reader, namesLen)
			if err != nil {
				utils.HandleError(e, "Error generating random number", err.Error())
				return err
			}
			secondIndex = int(secondBig.Int64())
			if secondIndex != firstIndex {
				break
			}
		}

		randomName := fmt.Sprintf("You are **%s%s**", names[firstIndex], names[secondIndex])

		_, err = e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetContent(randomName).
			Build())

		return err
	}
}
