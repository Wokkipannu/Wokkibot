package utils

import "strings"

// Returns a color integer from RGB values
func RGBToInteger(r, g, b int) int {
	return (r << 16) + (g << 8) + b
}

func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}
