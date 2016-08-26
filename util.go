package main

import "strings"

func ReadableTitle(title string) string {
	for _, prefix := range []string{
		"A", "The", "Der", "Le", "Les", "La", "Las", "El",
	} {
		suffix := ", " + prefix
		if !strings.HasSuffix(title, suffix) {
			continue
		}
		title = prefix + " " + strings.TrimSuffix(title, suffix)
	}
	return title
}

func PageTitle(title string) string {
	result := "Starring the Computer"
	if title != "" {
		result += " - " + title
	}
	return result
}
