package utils

import (
	"strings"
	"unicode"
)

func GenerateSlug(input string) string {
	input = strings.ToLower(input)
	var sb strings.Builder
	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			sb.WriteRune(r)
		} else if r == ' ' || r == '-' {
			sb.WriteRune('-')
		}
	}
	return strings.Trim(sb.String(), "-")
}
