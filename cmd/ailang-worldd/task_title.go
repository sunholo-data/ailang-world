package main

import "strings"

// taskTitle derives a short title for new starts; stored titles stay unchanged.
func taskTitle(data []byte) string {
	title := ""
	for _, line := range strings.Split(string(data), "\n") {
		title = strings.TrimSpace(line)
		if title != "" {
			break
		}
	}
	if strings.HasPrefix(title, "#") {
		title = strings.TrimSpace(strings.TrimLeft(title, "#"))
	} else {
		if len(title) >= 5 && strings.EqualFold(title[:5], "Task:") {
			title = strings.TrimSpace(title[5:])
		}
		for _, sep := range []string{". ", "; ", " ("} {
			if i := strings.Index(title, sep); i >= 0 {
				title = title[:i]
			}
		}
	}
	title = strings.Join(strings.Fields(title), " ")
	runes := []rune(title)
	if len(runes) > 72 {
		cut := 72
		if i := strings.LastIndex(string(runes[:72]), " "); i >= 0 {
			cut = len([]rune(string(runes[:72])[:i]))
		}
		title = string(runes[:cut]) + "…"
	}
	return title
}
