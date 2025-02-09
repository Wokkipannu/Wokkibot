package trivia

import (
	"regexp"
	"strings"
	"time"
	"unicode"
)

type AnswerValidator struct {
	correctAnswer string
}

func NewAnswerValidator(correctAnswer string) *AnswerValidator {
	return &AnswerValidator{
		correctAnswer: correctAnswer,
	}
}

func (v *AnswerValidator) ValidateAnswer(userAnswer string) bool {
	return v.validateDate(userAnswer) ||
		v.validateNumber(userAnswer) ||
		v.validateName(userAnswer) ||
		v.validateGeneral(userAnswer)
}

func (v *AnswerValidator) validateName(userAnswer string) bool {
	user := splitName(cleanString(userAnswer))
	correct := splitName(cleanString(v.correctAnswer))

	if !containsOnlyLetters(userAnswer) || !containsOnlyLetters(v.correctAnswer) {
		return false
	}

	if len(user) == 1 {
		for _, part := range correct {
			if strings.EqualFold(user[0], part) {
				return true
			}
		}
		return false
	}

	if len(user) != len(correct) {
		return false
	}

	matches := 0
	for i := range user {
		if strings.EqualFold(user[i], correct[i]) {
			matches++
		}
	}

	return matches == len(correct)
}

func (v *AnswerValidator) validateDate(userAnswer string) bool {
	userDate := parseDate(userAnswer)
	correctDate := parseDate(v.correctAnswer)

	if userDate == nil || correctDate == nil {
		return false
	}

	if userDate.Day() == 1 && userDate.Month() == 1 {
		return userDate.Year() == correctDate.Year()
	}

	if userDate.Day() == 1 {
		return userDate.Year() == correctDate.Year() &&
			userDate.Month() == correctDate.Month()
	}

	return userDate.Equal(*correctDate)
}

func (v *AnswerValidator) validateNumber(userAnswer string) bool {
	userNum := extractNumber(cleanString(userAnswer))
	correctNum := extractNumber(cleanString(v.correctAnswer))

	if userNum == "" || correctNum == "" {
		return false
	}

	return userNum == correctNum
}

func (v *AnswerValidator) validateGeneral(userAnswer string) bool {
	user := cleanString(userAnswer)
	correct := cleanString(v.correctAnswer)

	if len(correct) <= 5 {
		return strings.EqualFold(user, correct)
	}

	similarity := calculateSimilarity(user, correct)
	return similarity >= 0.85
}

func cleanString(s string) string {
	s = strings.Map(func(r rune) rune {
		if unicode.IsPunct(r) {
			return -1
		}
		return r
	}, s)
	return strings.TrimSpace(strings.ToLower(s))
}

func splitName(name string) []string {
	return strings.Fields(name)
}

func containsOnlyLetters(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

func parseDate(date string) *time.Time {
	formats := []string{
		"2006",
		"January 2006",
		"Jan 2006",
		"2/2006",
		"01/02/2006",
		"2006-01-02",
		"January 2, 2006",
		"Jan 2, 2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, date); err == nil {
			return &t
		}
	}

	re := regexp.MustCompile(`\b(19|20)\d{2}\b`)
	if year := re.FindString(date); year != "" {
		if t, err := time.Parse("2006", year); err == nil {
			return &t
		}
	}

	return nil
}

func extractNumber(s string) string {
	re := regexp.MustCompile(`[0-9]+`)
	return re.FindString(s)
}

func calculateSimilarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}

	longer := s1
	shorter := s2
	if len(s1) < len(s2) {
		longer = s2
		shorter = s1
	}

	longerLength := len(longer)
	if longerLength == 0 {
		return 1.0
	}

	return (float64(longerLength) - float64(editDistance(longer, shorter))) / float64(longerLength)
}

func editDistance(s1, s2 string) int {
	m := len(s1)
	n := len(s2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 0; i <= m; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s1[i-1] == s2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = 1 + min(dp[i-1][j], min(dp[i][j-1], dp[i-1][j-1]))
			}
		}
	}

	return dp[m][n]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
