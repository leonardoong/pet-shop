package slug

import (
	"fmt"
	"regexp"
	"strings"
)

var nonAlphaRegex = regexp.MustCompile(`[^a-z0-9\-]+`)
var multiHyphenRegex = regexp.MustCompile(`-+`)

func Make(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	s = nonAlphaRegex.ReplaceAllString(s, "")
	s = multiHyphenRegex.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

func Unique(s string, exists func(string) (bool, error)) (string, error) {
	base := Make(s)
	slug := base
	n := 1
	for {
		ok, err := exists(slug)
		if err != nil {
			return "", err
		}
		if !ok {
			return slug, nil
		}
		n++
		slug = Make(fmt.Sprintf("%s-%d", base, n))
	}
}
