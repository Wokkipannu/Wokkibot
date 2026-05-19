package user

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"
	"wokkibot/common"
	"wokkibot/utils"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
)

// TODO: move to config file eventually
const badgeEmojiGuildID snowflake.ID = 1204502217470513203

var badgeEmojiCache = struct {
	sync.Mutex
	loaded   bool
	mentions map[string]string
}{}

var UserCommand = discord.SlashCommandCreate{
	Name:        "user",
	Description: "Get information about a user",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionUser{
			Name:        "user",
			Description: "The user to get information about",
			Required:    false,
		},
	},
}

func HandleUser(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		data := e.SlashCommandInteractionData()

		var user discord.User

		if u, ok := data.OptUser("user"); ok {
			user = u
		} else {
			user = e.User()
		}

		// For some reason, user does not contain certain attributes, such as BannerURL or AccentColor, so we must fetch the user from the client rest
		fetchedUser, err := b.Client.Rest.GetUser(user.ID)
		if err != nil {
			slog.Info("Error fetching user from client")
		}
		if fetchedUser != nil {
			user = *fetchedUser
		}

		userBadges := formatUserBadges(b, user.PublicFlags)

		embed := discord.NewEmbed()
		embed = embed.WithAuthor(fmt.Sprintf("%v's profile", user.EffectiveName()), "", *user.AvatarURL())
		if user.Bot {
			embed = embed.AddField("Type", "Bot", true)
		} else {
			embed = embed.AddField("Type", "User", true)
		}
		embed = embed.AddField("Global name", user.EffectiveName(), true)
		embed = embed.AddField("Username", user.Username, true)
		if len(userBadges) > 0 {
			embed = embed.AddField("Badges", strings.Join(userBadges, " "), true)
		}
		if joinedAt := e.Member().JoinedAt; joinedAt != nil {
			embed = embed.AddField("Joined this server", fmt.Sprintf("%v (%v days ago)", joinedAt.Format("02.01.2006"), DaysSince(*joinedAt)), false)
		}
		embed = embed.AddField("Account created", fmt.Sprintf("%v (%v days ago)", user.CreatedAt().Format("02.01.2006"), DaysSince(user.CreatedAt())), false)

		embed = embed.WithThumbnail(user.EffectiveAvatarURL())

		if user.AccentColor != nil {
			embed = embed.WithColor(*user.AccentColor)
		}

		if user.BannerURL() != nil {
			formatOpt := utils.SetCDNOptions(discord.FileFormatPNG, discord.QueryValues{"size": 1024})
			embed = embed.WithImage(*user.BannerURL(formatOpt))
		}

		return e.CreateMessage(discord.NewMessageCreate().WithEmbeds(embed))
	}
}

func DaysSince(date time.Time) int {
	return int(time.Since(date).Hours() / 24)
}

func formatUserBadges(b *wokkibot.Wokkibot, flags discord.UserFlags) []string {
	emojiMentions := getBadgeEmojiMentions(b)
	badges := make([]string, 0, len(common.UserBadges))

	for _, badge := range common.UserBadges {
		if !flags.Has(badge.Flag) {
			continue
		}

		if badge.EmojiName != "" {
			if mention, ok := emojiMentions[badge.EmojiName]; ok {
				badges = append(badges, mention)
				continue
			}
		}

		badges = append(badges, badge.Name)
	}

	return badges
}

func getBadgeEmojiMentions(b *wokkibot.Wokkibot) map[string]string {
	badgeEmojiCache.Lock()
	defer badgeEmojiCache.Unlock()

	if badgeEmojiCache.loaded {
		return badgeEmojiCache.mentions
	}

	emojis, err := b.Client.Rest.GetEmojis(badgeEmojiGuildID)
	if err != nil {
		slog.Warn("failed to fetch user badge emojis", "guild_id", badgeEmojiGuildID, "error", err)
		return nil
	}

	mentions := make(map[string]string, len(emojis))
	for _, emoji := range emojis {
		mentions[emoji.Name] = emoji.Mention()
	}

	badgeEmojiCache.mentions = mentions
	badgeEmojiCache.loaded = true
	return mentions
}
