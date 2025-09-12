package tokenizer

import (
	"strings"
	"unicode"
)

type Token struct {
	Tokens []string
}

// Remove leading and tailing punctuation
func (t *Token) removePunctuation() {
	tok := make([]string, 0, len(t.Tokens))
	for _, token := range t.Tokens {
		token = cleanToken(token)
		if len(token) != 0 {
			tok = append(tok, token)
		}
	}
	t.Tokens = tok
}

func (t *Token) toLower() {
	for i := range t.Tokens {
		t.Tokens[i] = strings.ToLower(t.Tokens[i])
	}
}

func cleanToken(token string) string {
	start, end := 0, len(token)-1
	for i := 0; i < len(token); i++ {
		if unicode.IsPunct(rune(token[i])) {
			continue
		}
		start = i
		break
	}
	for i := len(token) - 1; i >= 0; i-- {
		if unicode.IsPunct(rune(token[i])) {
			continue
		}
		end = i
		break
	}

	var builder strings.Builder
	for i := start; i <= end; i++ {
		builder.WriteRune(rune(token[i]))
	}
	return builder.String()
}
