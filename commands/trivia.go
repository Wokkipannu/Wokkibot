package commands

import (
	"encoding/json"
	"html"
	"io"
	"log/slog"
	"math/rand"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"
	"wokkibot/utils"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"golang.org/x/net/context"
)

type Res struct {
	Code   int              `json:"response_code"`
	Trivia []TriviaQuestion `json:"results"`
}

type TriviaQuestion struct {
	Type             string   `json:"type"`
	Difficulty       string   `json:"difficulty"`
	Category         string   `json:"category"`
	Question         string   `json:"question"`
	CorrectAnswer    string   `json:"correct_answer"`
	IncorrectAnswers []string `json:"incorrect_answers"`
}

var triviaCommand = discord.SlashCommandCreate{
	Name:        "trivia",
	Description: "Play a trivia game",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "category",
			Description: "Category of the trivia",
			Required:    false,
			Choices: []discord.ApplicationCommandOptionChoiceString{
				{
					Name:  "General Knowledge",
					Value: "9",
				},
				{
					Name:  "Entertainment: Books",
					Value: "10",
				},
				{
					Name:  "Entertainment: Film",
					Value: "11",
				},
				{
					Name:  "Entertainment: Music",
					Value: "12",
				},
				{
					Name:  "Entertainment: Musicals & Theatres",
					Value: "13",
				},
				{
					Name:  "Entertainment: Television",
					Value: "14",
				},
				{
					Name:  "Entertainment: Video Games",
					Value: "15",
				},
				{
					Name:  "Entertainment: Board Games",
					Value: "16",
				},
				{
					Name:  "Entertainment: Japanese Anime & Manga",
					Value: "31",
				},
				{
					Name:  "Entertainment: Cartoon & Animations",
					Value: "32",
				},
				{
					Name:  "Entertainment: Comics",
					Value: "29",
				},
				{
					Name:  "Science & Nature",
					Value: "17",
				},
				{
					Name:  "Science: Computers",
					Value: "18",
				},
				{
					Name:  "Science: Mathematics",
					Value: "19",
				},
				{
					Name:  "Science: Gadgets",
					Value: "30",
				},
				{
					Name:  "Mythology",
					Value: "20",
				},
				{
					Name:  "Sports",
					Value: "21",
				},
				{
					Name:  "Geography",
					Value: "22",
				},
				{
					Name:  "History",
					Value: "23",
				},
				{
					Name:  "Politics",
					Value: "24",
				},
				{
					Name:  "Art",
					Value: "25",
				},
				{
					Name:  "Celebrities",
					Value: "26",
				},
				{
					Name:  "Animals",
					Value: "27",
				},
				{
					Name:  "Vehicles",
					Value: "28",
				},
			},
		},
		discord.ApplicationCommandOptionString{
			Name:        "difficulty",
			Description: "Difficulty of the trivia",
			Required:    false,
			Choices: []discord.ApplicationCommandOptionChoiceString{
				{
					Name:  "Easy",
					Value: "easy",
				},
				{
					Name:  "Medium",
					Value: "medium",
				},
				{
					Name:  "Hard",
					Value: "hard",
				},
			},
		},
		discord.ApplicationCommandOptionString{
			Name:        "type",
			Description: "Type of the trivia",
			Required:    false,
			Choices: []discord.ApplicationCommandOptionChoiceString{
				{
					Name:  "Multiple Choice (Default)",
					Value: "multiple",
				},
				{
					Name:  "True / False",
					Value: "boolean",
				},
			},
		},
	},
}

func HandleTrivia(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		data := e.SlashCommandInteractionData()

		t := b.Trivias.Get(*e.GuildID())

		if t.IsActive {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Trivia is already running. Wait for it to finish first.").Build())
		}

		apiEndpoint := "https://opentdb.com/api.php"
		queryParams := url.Values{}

		queryParams.Add("amount", "1")

		if category, ok := data.OptString("category"); ok {
			queryParams.Add("category", category)
		}

		if difficulty, ok := data.OptString("difficulty"); ok {
			queryParams.Add("difficulty", difficulty)
		}
		if type_, ok := data.OptString("type"); ok {
			queryParams.Add("type", type_)
		} else {
			queryParams.Add("type", "multiple")
		}

		res, err := b.Client.Rest().HTTPClient().Get(apiEndpoint + "?" + queryParams.Encode())
		if err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Error while fetching data").Build())
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Error while reading data").Build())
		}

		var triviaResponse Res
		err = json.Unmarshal(body, &triviaResponse)
		if err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Error while parsing data").Build())
		}

		if len(triviaResponse.Trivia) == 0 {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("API did not return any trivia. Try again later.").Build())
		}
		trivia := triviaResponse.Trivia[0]

		embed := discord.NewEmbedBuilder()
		embed.SetTitle("Trivia Question")
		embed.AddField("Difficulty", trivia.Difficulty, true)
		embed.AddField("Category", trivia.Category, true)
		embed.SetColor(utils.RGBToInteger(255, 215, 0))
		embed.SetDescription(html.UnescapeString(trivia.Question))
		embed.SetFooterText("Type your answers below. Time limit 60 seconds.")
		// embed.AddField("Correct answer", fmt.Sprintf("||%v||", trivia.CorrectAnswer), true)
		// embed.AddField("Answers", strings.Join(trivia.IncorrectAnswers, "\n"), true)

		go func(channel snowflake.ID) {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("Recovered from panic in goroutine", slog.Any("recover", r))
				}
			}()

			ec := wokkibot.NewEventCollector(b.Client, channel)
			defer func() {
				ec.Stop()
			}()

			ctx, clsCtx := context.WithTimeout(context.Background(), 60*time.Second)
			defer clsCtx()

			for {
				select {
				case <-ctx.Done():
					_, err := b.Client.Rest().CreateMessage(channel, discord.NewMessageCreateBuilder().SetContentf("Trivia timed out. The correct answer was %v", trivia.CorrectAnswer).Build())
					if err != nil {
						slog.Error("Error while sending timeout message", slog.Any("err", err))
					}
					t.SetStatus(false)
					return

				case messageEvent := <-ec.Events():
					if messageEvent == nil {
						slog.Warn("Received nil message event", slog.Any("channel", channel))
						continue
					}

					distance := levenshtein.DistanceForStrings([]rune(strings.ToLower(trivia.CorrectAnswer)), []rune(strings.ToLower(messageEvent.Message.Content)), levenshtein.DefaultOptions)

					threshold := utf8.RuneCountInString(trivia.CorrectAnswer) / 2

					if distance <= threshold {
						_, err := b.Client.Rest().CreateMessage(messageEvent.ChannelID, discord.NewMessageCreateBuilder().SetContentf("%v got it correct! The correct answer was %v", messageEvent.Message.Author.EffectiveName(), trivia.CorrectAnswer).SetMessageReference(messageEvent.Message.MessageReference).Build())
						if err != nil {
							slog.Error("Error while sending correct answer message", slog.Any("err", err))
						}
						t.SetStatus(false)
						return
					}

					if strings.ToLower(messageEvent.Message.Content) == "hint" {
						var options []string
						options = append(options, trivia.IncorrectAnswers...)
						options = append(options, trivia.CorrectAnswer)

						for i := range options {
							j := rand.Intn(i + 1)
							options[i], options[j] = html.UnescapeString(options[j]), html.UnescapeString(options[i])
						}

						hintEmbed := discord.NewEmbedBuilder()
						hintEmbed.SetTitle("Trivia Hint")
						hintEmbed.SetDescription(trivia.Question)
						hintEmbed.AddField("Choices", strings.Join(options, "\n"), true)

						_, err := b.Client.Rest().CreateMessage(messageEvent.ChannelID, discord.NewMessageCreateBuilder().SetEmbeds(hintEmbed.Build()).SetMessageReference(messageEvent.Message.MessageReference).Build())
						if err != nil {
							slog.Error("Error while sending hint", slog.Any("err", err))
						}
					}
				}
			}
		}(e.Channel().ID())

		t.SetStatus(true)
		return e.CreateMessage(discord.NewMessageCreateBuilder().SetEmbeds(embed.Build()).Build())
	}
}
