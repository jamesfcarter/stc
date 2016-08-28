package main

import (
	"database/sql/driver"
	"fmt"
	"html"
	"html/template"
	"regexp"
	"strconv"
	"strings"
)

type Markup string

var (
	MarkupRe = regexp.MustCompile("(FEAT|COMP|HELP)\\[[^\\]]*\\]")
)

func replace(s string, fullUrl bool) string {
	s = strings.TrimRight(s, "]")
	link := ""
	for prefix, l := range map[string]string{
		"FEAT[": "/feature.html?f=",
		"COMP[": "/computer.html?c=",
		"HELP[": "/help.html?h=",
	} {
		if strings.HasPrefix(s, prefix) {
			s = strings.TrimPrefix(s, prefix)
			link = l
			break
		}
	}
	s = strings.TrimSpace(s)
	spl := strings.SplitN(s, ",", 2)
	id, _ := strconv.Atoi(spl[0])
	s = strings.TrimSpace(spl[1])

	link += strconv.Itoa(id)
	if fullUrl {
		link = "http://starringthecomputer.com" + link
	}

	return fmt.Sprintf("<a href='%s'>%s</a>", link, html.EscapeString(s))
}

func (m Markup) Format() template.HTML {
	unformatted := string(m)

	return template.HTML(MarkupRe.ReplaceAllStringFunc(unformatted,
		func(s string) string {
			return replace(s, false)
		}))
}

func (m Markup) FormatFullUrl() template.HTML {
	unformatted := string(m)

	return template.HTML(MarkupRe.ReplaceAllStringFunc(unformatted,
		func(s string) string {
			return replace(s, true)
		}))
}

func (m *Markup) Scan(value interface{}) error {
	*m = ""
	if value != nil {
		*m = Markup(value.([]uint8))
	}
	return nil
}

func (m Markup) Value() (driver.Value, error) {
	return m, nil
}
