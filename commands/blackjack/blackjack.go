package blackjack

import (
	"fmt"
	"log/slog"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"wokkibot/utils"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
	"golang.org/x/net/context"
)

var suits = []string{"♠️", "♥️", "♦️", "♣️"}
var ranks = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

type Card struct {
	Suit  string
	Rank  string
	Value int
}

func (c Card) String() string {
	return fmt.Sprintf("%s%s", c.Rank, c.Suit)
}

type Hand struct {
	Cards []Card
}

func (h *Hand) AddCard(card Card) {
	h.Cards = append(h.Cards, card)
}

func (h *Hand) Value() int {
	value := 0
	aces := 0

	for _, card := range h.Cards {
		if card.Rank == "A" {
			aces++
			value += 11
		} else {
			value += card.Value
		}
	}

	for aces > 0 && value > 21 {
		value -= 10
		aces--
	}

	return value
}

func (h *Hand) String() string {
	cards := make([]string, len(h.Cards))
	for i, card := range h.Cards {
		cards[i] = card.String()
	}
	return strings.Join(cards, " ")
}

func (h *Hand) IsBusted() bool {
	return h.Value() > 21
}

func (h *Hand) IsBlackjack() bool {
	return len(h.Cards) == 2 && h.Value() == 21
}

type Player struct {
	User        discord.User
	DisplayName string
	Hand        Hand
	Standing    bool
	Busted      bool
}

type Dealer struct {
	User        *discord.User
	DisplayName string
	Hand        Hand
}

func (d *Dealer) Name() string {
	if d.DisplayName != "" {
		return d.DisplayName
	}
	if d.User != nil {
		return d.User.EffectiveName()
	}
	return "Dealer"
}

func (d *Dealer) FormattedName() string {
	if d.DisplayName != "" {
		return fmt.Sprintf("🎰 %s (Dealer)", d.DisplayName)
	}
	if d.User != nil {
		return fmt.Sprintf("🎰 %s (Dealer)", d.User.EffectiveName())
	}
	return "🎰 Dealer"
}

type Deck struct {
	Cards []Card
}

func NewDeck() *Deck {
	deck := &Deck{}
	for _, suit := range suits {
		for _, rank := range ranks {
			value := 0
			switch rank {
			case "A":
				value = 11
			case "J", "Q", "K":
				value = 10
			default:
				fmt.Sscanf(rank, "%d", &value)
			}
			deck.Cards = append(deck.Cards, Card{Suit: suit, Rank: rank, Value: value})
		}
	}
	return deck
}

func (d *Deck) Shuffle() {
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

func (d *Deck) Draw() Card {
	card := d.Cards[0]
	d.Cards = d.Cards[1:]
	return card
}

var BlackjackCommand = discord.SlashCommandCreate{
	Name:        "blackjack",
	Description: "Play a game of blackjack",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "players",
			Description: "Mention the players (e.g. @Player1 @Player2 @Player3)",
			Required:    true,
		},
		discord.ApplicationCommandOptionString{
			Name:        "dealer",
			Description: "Mention a user to be the dealer (bot is dealer by default)",
			Required:    false,
		},
	},
}

func HandleBlackjack(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		if err := e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, nil); err != nil {
			return err
		}

		game := b.Blackjacks.Get(*e.GuildID())

		if game.IsActive {
			_, err := e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
				SetContent("A blackjack game is already running in this server. Wait for it to finish first.").
				Build())
			return err
		}

		data := e.SlashCommandInteractionData()
		playersInput := data.String("players")

		mentionRegex := regexp.MustCompile(`<@!?(\d+)>`)
		matches := mentionRegex.FindAllStringSubmatch(playersInput, -1)

		if len(matches) == 0 {
			_, err := e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
				SetContent("Please mention at least one player (e.g. /blackjack players:@Player1 @Player2)").
				Build())
			return err
		}

		var players []Player
		seenIDs := make(map[snowflake.ID]bool)
		guildID := *e.GuildID()
		for _, match := range matches {
			userID, err := snowflake.Parse(match[1])
			if err != nil {
				continue
			}

			if seenIDs[userID] {
				continue
			}
			seenIDs[userID] = true

			user, err := b.Client.Rest().GetUser(userID)
			if err != nil {
				slog.Warn("Failed to fetch user", slog.Any("userID", userID), slog.Any("err", err))
				continue
			}

			if !user.Bot {
				displayName := user.EffectiveName()
				member, err := b.Client.Rest().GetMember(guildID, userID)
				if err == nil && member.Nick != nil && *member.Nick != "" {
					displayName = *member.Nick
				}

				players = append(players, Player{
					User:        *user,
					DisplayName: displayName,
					Hand:        Hand{},
				})
			}
		}

		if len(players) == 0 {
			_, err := e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
				SetContent("No valid players found. Bots cannot play!").
				Build())
			return err
		}

		if len(players) > 7 {
			_, err := e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
				SetContent("Maximum 7 players allowed in a blackjack game.").
				Build())
			return err
		}

		var dealer Dealer
		if dealerInput, ok := data.OptString("dealer"); ok && dealerInput != "" {
			dealerMatches := mentionRegex.FindAllStringSubmatch(dealerInput, 1)
			if len(dealerMatches) > 0 {
				dealerID, err := snowflake.Parse(dealerMatches[0][1])
				if err == nil {
					dealerUser, err := b.Client.Rest().GetUser(dealerID)
					if err == nil && !dealerUser.Bot {
						dealer.User = dealerUser
						dealer.DisplayName = dealerUser.EffectiveName()
						member, err := b.Client.Rest().GetMember(guildID, dealerID)
						if err == nil && member.Nick != nil && *member.Nick != "" {
							dealer.DisplayName = *member.Nick
						}
					}
				}
			}
		}

		numDecks := (len(players) / 3) + 1
		var allCards []Card
		for i := 0; i < numDecks; i++ {
			deck := NewDeck()
			allCards = append(allCards, deck.Cards...)
		}
		deck := &Deck{Cards: allCards}
		deck.Shuffle()

		for i := range players {
			players[i].Hand.AddCard(deck.Draw())
			players[i].Hand.AddCard(deck.Draw())
		}
		dealer.Hand.AddCard(deck.Draw())
		dealer.Hand.AddCard(deck.Draw())

		embed := buildGameEmbed(players, dealer, -1, true)

		_, err := e.UpdateInteractionResponse(discord.NewMessageUpdateBuilder().
			SetContent("").
			SetEmbeds(embed.Build()).
			Build())
		if err != nil {
			return err
		}

		utils.UpdateStatistics("blackjack_games_played")

		go runBlackjackGame(b, e, players, dealer, deck)

		game.SetStatus(true)
		_ = playersInput
		return nil
	}
}

func runBlackjackGame(b *wokkibot.Wokkibot, e *handler.CommandEvent, players []Player, dealer Dealer, deck *Deck) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Recovered from panic in blackjack goroutine", slog.Any("recover", r))
		}
		b.Blackjacks.Get(*e.GuildID()).SetStatus(false)
	}()

	channelID := e.Channel().ID()

	for playerIdx := range players {
		if players[playerIdx].Hand.IsBlackjack() {
			embed := buildGameEmbed(players, dealer, playerIdx, true)
			embed.SetFooter(fmt.Sprintf("🎰 %s has BLACKJACK!", players[playerIdx].DisplayName), "")
			_, _ = b.Client.Rest().CreateMessage(channelID, discord.NewMessageCreateBuilder().SetEmbeds(embed.Build()).Build())
			time.Sleep(2 * time.Second)
			continue
		}

		for !players[playerIdx].Standing && !players[playerIdx].Busted {
			embed := buildGameEmbed(players, dealer, playerIdx, true)
			embed.SetFooter(fmt.Sprintf("🎯 %s's turn! Type 'hit' or 'stand' (60s timeout)", players[playerIdx].DisplayName), "")

			_, _ = b.Client.Rest().CreateMessage(channelID, discord.NewMessageCreateBuilder().SetEmbeds(embed.Build()).Build())

			action, timedOut := waitForPlayerAction(e, players[playerIdx].User.ID, channelID)

			if timedOut {
				players[playerIdx].Standing = true
				_, _ = b.Client.Rest().CreateMessage(channelID, discord.NewMessageCreateBuilder().
					SetContent(fmt.Sprintf("⏰ %s took too long! Auto-standing.", players[playerIdx].User.Mention())).
					Build())
				continue
			}

			switch action {
			case "hit":
				players[playerIdx].Hand.AddCard(deck.Draw())
				if players[playerIdx].Hand.IsBusted() {
					players[playerIdx].Busted = true
					embed := buildGameEmbed(players, dealer, playerIdx, true)
					embed.SetFooter(fmt.Sprintf("💥 %s busted with %d!", players[playerIdx].DisplayName, players[playerIdx].Hand.Value()), "")
					_, _ = b.Client.Rest().CreateMessage(channelID, discord.NewMessageCreateBuilder().SetEmbeds(embed.Build()).Build())
				}
			case "stand":
				players[playerIdx].Standing = true
			}
		}
	}

	dealerName := dealer.Name()
	_, _ = b.Client.Rest().CreateMessage(channelID, discord.NewMessageCreateBuilder().
		SetContent(fmt.Sprintf("🎰 **%s's turn!**", dealerName)).
		Build())

	time.Sleep(1 * time.Second)

	embed := buildGameEmbed(players, dealer, -1, false)
	embed.SetFooter(fmt.Sprintf("%s shows: %s (Value: %d)", dealerName, dealer.Hand.String(), dealer.Hand.Value()), "")
	_, _ = b.Client.Rest().CreateMessage(channelID, discord.NewMessageCreateBuilder().SetEmbeds(embed.Build()).Build())

	allBusted := true
	for _, p := range players {
		if !p.Busted {
			allBusted = false
			break
		}
	}

	if !allBusted {
		for dealer.Hand.Value() < 17 {
			time.Sleep(1500 * time.Millisecond)
			dealer.Hand.AddCard(deck.Draw())

			embed := buildGameEmbed(players, dealer, -1, false)
			if dealer.Hand.IsBusted() {
				embed.SetFooter(fmt.Sprintf("💥 %s busted with %d!", dealerName, dealer.Hand.Value()), "")
			} else {
				embed.SetFooter(fmt.Sprintf("%s draws... (Value: %d)", dealerName, dealer.Hand.Value()), "")
			}
			_, _ = b.Client.Rest().CreateMessage(channelID, discord.NewMessageCreateBuilder().SetEmbeds(embed.Build()).Build())
		}
	}

	time.Sleep(1 * time.Second)
	resultsEmbed := buildResultsEmbed(players, dealer)
	_, _ = b.Client.Rest().CreateMessage(channelID, discord.NewMessageCreateBuilder().SetEmbeds(resultsEmbed.Build()).Build())
}

func waitForPlayerAction(e *handler.CommandEvent, playerID snowflake.ID, channelID snowflake.ID) (string, bool) {
	ch, cls := bot.NewEventCollector(e.Client(), func(event *events.MessageCreate) bool {
		if event.Message.Author.Bot || event.ChannelID != channelID {
			return false
		}
		if event.Message.Author.ID != playerID {
			return false
		}
		content := strings.ToLower(strings.TrimSpace(event.Message.Content))
		return content == "hit" || content == "stand"
	})
	defer cls()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return "", true
	case msg := <-ch:
		if msg == nil {
			return "", true
		}
		return strings.ToLower(strings.TrimSpace(msg.Message.Content)), false
	}
}

func buildGameEmbed(players []Player, dealer Dealer, currentPlayerIdx int, hideDealer bool) *discord.EmbedBuilder {
	embed := discord.NewEmbedBuilder()
	embed.SetTitle("🃏 Blackjack")
	embed.SetColor(utils.RGBToInteger(46, 139, 87))

	if hideDealer && len(dealer.Hand.Cards) > 0 {
		embed.AddField(dealer.FormattedName(), fmt.Sprintf("%s 🂠", dealer.Hand.Cards[0].String()), false)
	} else {
		embed.AddField(dealer.FormattedName(), fmt.Sprintf("%s (Value: %d)", dealer.Hand.String(), dealer.Hand.Value()), false)
	}

	for i, player := range players {
		status := ""
		if player.Busted {
			status = " 💥 BUSTED"
		} else if player.Standing {
			status = " ✋ STANDING"
		} else if player.Hand.IsBlackjack() {
			status = " 🎰 BLACKJACK!"
		}

		indicator := ""
		if i == currentPlayerIdx {
			indicator = "▶ "
		}

		embed.AddField(
			fmt.Sprintf("%s%s", indicator, player.DisplayName),
			fmt.Sprintf("%s (Value: %d)%s", player.Hand.String(), player.Hand.Value(), status),
			true,
		)
	}

	return embed
}

func buildResultsEmbed(players []Player, dealer Dealer) *discord.EmbedBuilder {
	embed := discord.NewEmbedBuilder()
	embed.SetTitle("🏆 Blackjack Results")
	embed.SetColor(utils.RGBToInteger(255, 215, 0))

	dealerValue := dealer.Hand.Value()
	dealerBusted := dealer.Hand.IsBusted()
	dealerBlackjack := dealer.Hand.IsBlackjack()

	embed.AddField(dealer.FormattedName(), fmt.Sprintf("%s (Value: %d)", dealer.Hand.String(), dealerValue), false)

	var results []string
	for _, player := range players {
		playerValue := player.Hand.Value()
		result := ""

		if player.Busted {
			result = "❌ LOSE (Busted)"
		} else if dealerBusted {
			result = fmt.Sprintf("✅ WIN (%s busted)", dealer.Name())
		} else if player.Hand.IsBlackjack() && !dealerBlackjack {
			result = "💲🤑🎰🤑💲 WIN (Blackjack)"
		} else if dealerBlackjack && !player.Hand.IsBlackjack() {
			result = fmt.Sprintf("❌ LOSE (%s blackjack)", dealer.Name())
		} else if playerValue > dealerValue {
			result = "✅ WIN"
		} else if playerValue < dealerValue {
			result = "❌ LOSE"
		} else {
			result = "🤝 PUSH (Tie)"
		}

		results = append(results, fmt.Sprintf("**%s**: %s (%d) - %s",
			player.DisplayName,
			player.Hand.String(),
			playerValue,
			result,
		))
	}

	embed.SetDescription(strings.Join(results, "\n"))
	return embed
}
