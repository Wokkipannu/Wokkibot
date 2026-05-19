package common

import "github.com/disgoorg/disgo/discord"

type UserBadge struct {
	Flag      discord.UserFlags
	Name      string
	EmojiName string
}

var UserBadges = []UserBadge{
	{Flag: discord.UserFlagDiscordEmployee, Name: "Discord employee", EmojiName: "Discordstaff"},
	{Flag: discord.UserFlagPartneredServerOwner, Name: "Partnered server owner"},
	{Flag: discord.UserFlagHypeSquadEvents, Name: "Hype squad events", EmojiName: "Hypesquad_Events_Badge"},
	{Flag: discord.UserFlagBugHunterLevel1, Name: "Bug hunter level 1", EmojiName: "Bug_hunter_badge"},
	{Flag: discord.UserFlagHouseBravery, Name: "House bravery", EmojiName: "Hypesquad_bravery_badge"},
	{Flag: discord.UserFlagHouseBrilliance, Name: "House brilliance", EmojiName: "Hypesquad_brilliance_badge"},
	{Flag: discord.UserFlagHouseBalance, Name: "House balance", EmojiName: "Hypesquad_balance_badge"},
	{Flag: discord.UserFlagEarlySupporter, Name: "Early supporter", EmojiName: "Early_supporter_badge"},
	{Flag: discord.UserFlagTeamUser, Name: "Team user"},
	{Flag: discord.UserFlagBugHunterLevel2, Name: "Bug hunter level 2", EmojiName: "Bug_buster_badge"},
	{Flag: discord.UserFlagVerifiedBot, Name: "Verified bot"},
	{Flag: discord.UserFlagEarlyVerifiedBotDeveloper, Name: "Early verified bot developer", EmojiName: "Verified_developer_badge"},
	{Flag: discord.UserFlagDiscordCertifiedModerator, Name: "Discord certified moderator", EmojiName: "Moderator_Programs_Alumni"},
}
