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

	"github.com/xrash/smetrics"
	"golang.org/x/text/unicode/norm"
)

// Returns a color integer from RGB values
func RGBToInteger(r, g, b int) int {
	return (r << 16) + (g << 8) + b
}

// Capitalize the first letter of a string
func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func IsNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func CleanNumericAnswer(input string) string {
	re := regexp.MustCompile(`[^\d]`)
	cleaned := re.ReplaceAllString(input, "")
	return cleaned
}

func StringMatch(a, b string) bool {
	a = RemoveDiacritics(a)
	b = RemoveDiacritics(b)

	distance := smetrics.WagnerFischer(a, b, 1, 1, 2)
	threshold := 2

	longest := float64(len(a))
	if len(b) > len(a) {
		longest = float64(len(b))
	}
	similarityRatio := (longest - float64(distance)) / longest

	lengthDifference := float64(len(a)) / float64(len(b))
	return distance <= threshold && similarityRatio > 0.8 && lengthDifference > 0.75 && lengthDifference < 1.25
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
