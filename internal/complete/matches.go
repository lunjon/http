package complete

import "strings"

func Matches(s string, items []string) []string {
	l := strings.ToLower(s)
	matches := []string{}

	for _, item := range items {
		i := strings.ToLower(item)
		if strings.Contains(i, l) {
			matches = append(matches, item)
		}
	}

	return matches
}
