package wokkibot

import (
	"github.com/disgoorg/snowflake/v2"
)

type TriviaQuestion struct {
	Type             string   `json:"type"`
	Difficulty       string   `json:"difficulty"`
	Category         string   `json:"category"`
	Question         string   `json:"question"`
	CorrectAnswer    string   `json:"correct_answer"`
	IncorrectAnswers []string `json:"incorrect_answers"`
}

type Trivia struct {
	IsActive bool
}

func (t *Trivia) SetStatus(status bool) {
	t.IsActive = status
}

type TriviaManager struct {
	trivias map[snowflake.ID]*Trivia
}

func (t *TriviaManager) Get(guildID snowflake.ID) *Trivia {
	trivia, ok := t.trivias[guildID]
	if !ok {
		trivia = &Trivia{
			IsActive: false,
		}
		t.trivias[guildID] = trivia
	}
	return trivia
}
