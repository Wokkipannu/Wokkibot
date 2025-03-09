package utils

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"
	"unicode"

	"math/rand"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"golang.org/x/text/unicode/norm"
)

// Capitalize the first letter of a string
func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func ExtractYear(dateStr string) (string, error) {
	dateFormats := []string{
		"January 2, 2006", "Jan 2, 2006", "2 January 2006", "2 Jan 2006",
		"02/01/2006", "01/02/2006", "2006-01-02", "02-01-2006",
		"2 Jan 06", "2006", "06",
	}

	for _, format := range dateFormats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return strconv.Itoa(t.Year()), nil
		}
	}

	re := regexp.MustCompile(`\b(20\d{2}|19\d{2})\b`)
	match := re.FindString(dateStr)
	if match != "" {
		return match, nil
	}

	return "", fmt.Errorf("no valid year found")
}

// Dump goroutines for debugging
func DumpGoroutines() {
	timestamp := time.Now().Format("2006-01-02T15-04-05")
	filename := fmt.Sprintf("goroutine_dump_%v.txt", timestamp)
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal("could not create goroutine dump file:", err)
	}
	defer f.Close()
	pprof.Lookup("goroutine").WriteTo(f, 1)
}

func RemoveDiacritics(input string) string {
	t := norm.NFD.String(input)
	result := strings.Map(func(r rune) rune {
		if unicode.Is(unicode.Mn, r) {
			return -1
		}
		return r
	}, t)
	return norm.NFC.String(result)
}

func GenerateRandomName(length int) string {
	const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = characters[seed.Intn(len(characters))]
	}
	return string(b)
}

func HandleError(e *handler.CommandEvent, message string, errorMessage string) error {
	_, err := e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
		SetEmbeds(discord.NewEmbedBuilder().
			SetTitle(message).
			SetDescription(errorMessage).
			SetColor(RGBToInteger(255, 0, 0)).
			Build()).
		SetContent("").
		SetFlags(discord.MessageFlagEphemeral).
		Build())

	if err != nil {
		return err
	}

	return nil
}

func CalculateMaximumFileSizeForGuild(guild discord.Guild) int {
	if guild.PremiumTier == discord.PremiumTier2 {
		return 50
	} else if guild.PremiumTier == discord.PremiumTier3 {
		return 100
	} else {
		return 10
	}
}
