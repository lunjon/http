package complete

import "strings"

type match struct {
	s     string
	lower string
}

func Complete(text string, items []string) (string, []string) {
	matches := []string{}
	textLower := strings.ToLower(text)

	for _, item := range items {
		itemLower := strings.ToLower(item)
		if strings.HasPrefix(itemLower, textLower) {
			matches = append(matches, item)
		}
	}

	if len(matches) == 0 {
		return text, matches
	}

	prefix := commonPrefix(matches)
	return prefix, matches
}

func commonPrefix(items []string) string {
	prefix := ""
	if len(items) == 0 {
		return prefix
	}

	index := 0
	done := false
	for !done {
		if index == len(items[0]) {
			break
		}

		ch := items[0][index]
		for _, item := range items[1:] {
			if item[index] != ch {
				done = true
				break
			}
		}
		if done {
			break
		}

		index++
		prefix += string(ch)
	}

	return prefix
}
