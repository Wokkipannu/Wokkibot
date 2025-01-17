package commands

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log/slog"
	"math/rand"
	"net/url"
	"strings"
	"time"
	"wokkibot/database"
	"wokkibot/utils"
	"wokkibot/wokkibot"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
	"golang.org/x/net/context"
)

type Res struct {
	Code   int              `json:"response_code"`
	Trivia []TriviaQuestion `json:"results"`
}

type TokenRes struct {
	Code    int    `json:"response_code"`
	Message string `json:"response_message"`
	Token   string `json:"token"`
}

type TriviaQuestion struct {
	Type             string   `json:"type"`
	Difficulty       string   `json:"difficulty"`
	Category         string   `json:"category"`
	Question         string   `json:"question"`
	CorrectAnswer    string   `json:"correct_answer"`
	IncorrectAnswers []string `json:"incorrect_answers"`
}

type Answers struct {
	ID           snowflake.ID
	AnswersCount int
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

var db *sql.DB
var triviaToken string

func HandleTrivia(b *wokkibot.Wokkibot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		t := b.Trivias.Get(*e.GuildID())

		if t.IsActive {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Trivia is already running. Wait for it to finish first.").Build())
		}

		db = database.GetDB()

		err := db.QueryRow("SELECT trivia_token FROM guilds WHERE id = ?", *e.GuildID()).Scan(&triviaToken)
		if err != nil {
			FetchToken(e, b)
		}

		_trivia, err := FetchTrivia(e, b)
		if err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Error while fetching trivia. Maybe the API is down?").Build())
		}

		if len(_trivia) == 0 {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("API did not return any trivia. Try again later.").Build())
		}
		trivia := _trivia[0]

		var options []string
		options = append(options, trivia.IncorrectAnswers...)
		options = append(options, trivia.CorrectAnswer)

		options = ShuffleOptions(options)

		embed := discord.NewEmbedBuilder()
		embed.SetTitle("Trivia Question")

		if strings.Contains(strings.ToLower(trivia.Question), "which") {
			embed.AddField("Choices", strings.Join(options, "\n"), false)
		}

		embed.AddField("Difficulty", trivia.Difficulty, true)
		embed.AddField("Category", html.UnescapeString(trivia.Category), true)
		embed.SetColor(utils.RGBToInteger(255, 215, 0))
		embed.SetDescription(html.UnescapeString(trivia.Question))
		embed.SetFooterText("Type your answers below. Time limit 60 seconds. You can type hint or skip if you are stuck.")
		// embed.AddField("Correct answer", fmt.Sprintf("||%v||", trivia.CorrectAnswer), true)
		// embed.AddField("Answers", strings.Join(trivia.IncorrectAnswers, "\n"), true)

		go func(channel snowflake.ID, options []string) {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("Recovered from panic in goroutine", slog.Any("recover", r))
				}
			}()

			ch, cls := bot.NewEventCollector(e.Client(), func(event *events.MessageCreate) bool {
				return !event.Message.Author.Bot && event.ChannelID == e.Channel().ID()
			})

			ctx, clsCtx := context.WithTimeout(context.Background(), 60*time.Second)
			defer clsCtx()

			var answers []Answers
			answersIndexMap := make(map[snowflake.ID]int)

			addOrUpdateUser := func(id snowflake.ID, answerCount int) {
				if index, exists := answersIndexMap[id]; exists {
					answers[index].AnswersCount = answers[index].AnswersCount + answerCount
				} else {
					user := Answers{ID: id, AnswersCount: answerCount}
					answers = append(answers, user)
					answersIndexMap[id] = len(answers) - 1
				}
			}

			getUserByID := func(id snowflake.ID) *Answers {
				if index, exists := answersIndexMap[id]; exists {
					return &answers[index]
				}
				return nil
			}

			for {
				select {
				case <-ctx.Done():
					timeoutEmbed := discord.NewEmbedBuilder()
					timeoutEmbed.SetTitle("Trivia ended")
					timeoutEmbed.SetDescription("No one guessed in time.")
					timeoutEmbed.AddField("Question", html.UnescapeString(trivia.Question), true)
					timeoutEmbed.AddField("Correct answer", html.UnescapeString(trivia.CorrectAnswer), true)
					timeoutEmbed.SetColor(utils.RGBToInteger(215, 0, 0))

					_, err := b.Client.Rest().CreateMessage(channel, discord.NewMessageCreateBuilder().SetEmbeds(timeoutEmbed.Build()).Build())
					if err != nil {
						slog.Error("Error while sending timeout message", slog.Any("err", err))
					}
					t.SetStatus(false)
					cls()
					return

				case messageEvent := <-ch:
					if messageEvent == nil {
						slog.Warn("Received nil message event", slog.Any("channel", channel))
						continue
					}

					if strings.ToLower(messageEvent.Message.Content) == "hint" {
						hintEmbed := discord.NewEmbedBuilder()
						hintEmbed.SetTitle("Trivia Hint")
						hintEmbed.SetDescription(html.UnescapeString(trivia.Question))
						hintEmbed.AddField("Choices", strings.Join(options, "\n"), true)
						hintEmbed.SetColor(utils.RGBToInteger(255, 215, 0))

						_, err := b.Client.Rest().CreateMessage(messageEvent.ChannelID, discord.NewMessageCreateBuilder().SetEmbeds(hintEmbed.Build()).SetMessageReferenceByID(messageEvent.Message.ID).Build())
						if err != nil {
							slog.Error("Error while sending hint", slog.Any("err", err))
						}
					} else if strings.ToLower(messageEvent.Message.Content) == "skip" {
						skipEmbed := discord.NewEmbedBuilder()
						skipEmbed.SetTitle("Trivia ended")
						skipEmbed.SetDescription("Trivia was skipped")
						skipEmbed.AddField("Question", html.UnescapeString(trivia.Question), true)
						skipEmbed.AddField("Correct answer", html.UnescapeString(trivia.CorrectAnswer), true)
						skipEmbed.SetColor(utils.RGBToInteger(215, 0, 0))

						_, err := b.Client.Rest().CreateMessage(messageEvent.ChannelID, discord.NewMessageCreateBuilder().SetEmbeds(skipEmbed.Build()).SetMessageReferenceByID(messageEvent.Message.ID).Build())
						if err != nil {
							slog.Error("Error while sending skip message", slog.Any("err", err))
						}
						t.SetStatus(false)
						cls()
						return
					} else {
						addOrUpdateUser(messageEvent.Message.Author.ID, 1)
					}

					if ValidateTriviaAnswer(messageEvent.Message.Content, html.UnescapeString(trivia.CorrectAnswer)) {
						a := getUserByID(messageEvent.Message.Author.ID)
						correctEmbed := discord.NewEmbedBuilder()
						correctEmbed.SetTitle("Trivia ended")
						if a.AnswersCount == 1 {
							correctEmbed.SetDescriptionf("%v got it correct in first try!", messageEvent.Message.Author.EffectiveName())
						} else {
							correctEmbed.SetDescriptionf("%v got it correct after %v answers!", messageEvent.Message.Author.EffectiveName(), a.AnswersCount)
						}
						correctEmbed.AddField("Question", html.UnescapeString(trivia.Question), true)
						correctEmbed.AddField("Correct answer", html.UnescapeString(trivia.CorrectAnswer), true)
						correctEmbed.SetColor(utils.RGBToInteger(0, 215, 0))

						_, err := b.Client.Rest().CreateMessage(messageEvent.ChannelID, discord.NewMessageCreateBuilder().SetEmbeds(correctEmbed.Build()).SetMessageReferenceByID(messageEvent.Message.ID).Build())
						if err != nil {
							slog.Error("Error while sending correct answer message", slog.Any("err", err))
						}
						t.SetStatus(false)
						cls()
						return
					}
				}
			}
		}(e.Channel().ID(), options)

		t.SetStatus(true)
		return e.CreateMessage(discord.NewMessageCreateBuilder().SetEmbeds(embed.Build()).Build())
	}
}

func ShuffleOptions(options []string) []string {
	for i := range options {
		j := rand.Intn(i + 1)
		options[i], options[j] = html.UnescapeString(options[j]), html.UnescapeString(options[i])
	}
	return options
}

func ValidateTriviaAnswer(answer, correct string) bool {
	if utils.IsNumeric(answer) {
		cleanedUserAnswer := utils.CleanNumericAnswer(answer)
		cleanedCorrectAnswer := utils.CleanNumericAnswer(correct)
		return cleanedUserAnswer == cleanedCorrectAnswer
	}

	if _, err := utils.ExtractYear(answer); err == nil {
		cleanedUserAnswer := utils.CleanNumericAnswer(answer)
		return cleanedUserAnswer == utils.CleanNumericAnswer(correct)
	}

	return utils.StringMatch(strings.ToLower(answer), strings.ToLower(correct))
}

func FetchToken(e *handler.CommandEvent, b *wokkibot.Wokkibot) error {
	res, err := b.Client.Rest().HTTPClient().Get("https://opentdb.com/api_token.php?command=request")
	if err != nil {
		slog.Error("Error while getting trivia token", slog.Any("err", err))
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Error while reading trivia token", slog.Any("err", err))
	}

	var tokenResponse TokenRes
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		slog.Error("Error while parsing trivia token", slog.Any("err", err))
	}

	result, err := db.Exec("UPDATE guilds SET trivia_token = ? WHERE id = ?", tokenResponse.Token, *e.GuildID())
	if err != nil {
		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Failed to update trivia token").Build())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Failed to update trivia token").Build())
	}

	if rowsAffected == 0 {
		_, err = db.Exec("INSERT INTO guilds (id, trivia_token) VALUES (?, ?)", *e.GuildID(), tokenResponse.Token)
		if err != nil {
			return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Failed to update trivia token").Build())
		}
	}

	triviaToken = tokenResponse.Token

	return nil
}

var tries = 0

func FetchTrivia(e *handler.CommandEvent, b *wokkibot.Wokkibot) ([]TriviaQuestion, error) {
	data := e.SlashCommandInteractionData()

	apiEndpoint := "https://opentdb.com/api.php"
	queryParams := url.Values{}

	queryParams.Add("amount", "1")
	queryParams.Add("token", triviaToken)

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
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var triviaResponse Res
	err = json.Unmarshal(body, &triviaResponse)
	if err != nil {
		return nil, err
	}

	switch triviaResponse.Code {
	case 3:
		FetchToken(e, b)
		tries++
		if tries > 3 {
			return nil, fmt.Errorf("token expired")
		}
		return FetchTrivia(e, b)
	case 5:
		return nil, fmt.Errorf("too many requests in a short period of time")
	}

	return triviaResponse.Trivia, nil
}
