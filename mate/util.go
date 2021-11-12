package mate

import (
	"strings"
)

func parseTag(tag string) (string, map[string]bool) {
	var (
		tagValue string
		tagFlags map[string]bool
	)

	parts := strings.Split(tag, ",")
	if len(parts) > 1 {
		tagFlags = make(map[string]bool)
		for _, tagFlag := range parts[1:] {
			tagFlags[tagFlag] = true
		}
	}

	tagValue = parts[0]

	return tagValue, tagFlags
}
