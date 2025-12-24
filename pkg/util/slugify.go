package util

import (
	"regexp"
	"strings"
	"unicode"
)

func Slugify(input string) string {
	input = strings.ToLower(input)

	input = strings.TrimSpace(input)

	reSpace := regexp.MustCompile(`[\s\p{Zs}]+`)
	input = reSpace.ReplaceAllString(input, "-")

	var builder strings.Builder
	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			builder.WriteRune(r)
		}
	}

	reDash := regexp.MustCompile(`-+`)
	slug := reDash.ReplaceAllString(builder.String(), "-")

	slug = strings.Trim(slug, "-")

	return slug
}
