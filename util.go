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

func IndexChar(s string) string {
	var first string
	for _, c := range s {
		first = string(c)
		break
	}
	return strings.ToUpper(first)
}

func NonBroken(s string) string {
	var result []rune
	for _, r := range s {
		if r == '\u0020' {
			result = append(result, '\u00A0')
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}
