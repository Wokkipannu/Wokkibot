package utils

import (
	"strings"

	"github.com/disgoorg/disgo/discord"
)

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

var UserFlags = map[discord.UserFlags]string{
	discord.UserFlagActiveDeveloper:           "Active developer",
	discord.UserFlagBugHunterLevel1:           "Bug hunter level 1",
	discord.UserFlagBugHunterLevel2:           "Bug hunter level 2",
	discord.UserFlagDiscordCertifiedModerator: "Discord certified moderator",
	discord.UserFlagDiscordEmployee:           "Discord employee",
	discord.UserFlagEarlySupporter:            "Early supporter",
	discord.UserFlagEarlyVerifiedBotDeveloper: "Early verified bot developer",
	discord.UserFlagHouseBalance:              "House balance",
	discord.UserFlagHouseBravery:              "House bravery",
	discord.UserFlagHouseBrilliance:           "House brilliance",
	discord.UserFlagHypeSquadEvents:           "Hype squad events",
	discord.UserFlagPartneredServerOwner:      "Partnered server owner",
	discord.UserFlagTeamUser:                  "Team user",
	discord.UserFlagVerifiedBot:               "Verified bot",
}
