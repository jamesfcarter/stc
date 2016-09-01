package main

import (
	"log"
	"net/http"
	"os/exec"
	"path"
	"strings"
	"unicode"
)

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
	var first rune
	for _, c := range s {
		first = c
		break
	}
	if unicode.IsDigit(first) {
		return "0"
	}
	return strings.ToUpper(string(first))
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

func SimpleForm(r *http.Request) map[string]string {
	form := map[string]string{}

	r.ParseForm()
	for key, val := range r.Form {
		if len(val) < 1 {
			continue
		}
		form[key] = val[0]
	}
	return form
}

func SendEmail(subject, body string) {
	cmd := exec.Command("Mail", "-s", subject, "james@jfc.org.uk")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Printf("SendEmail1: %v", err)
		return
	}

	err = cmd.Start()
	if err != nil {
		log.Printf("SendEmail2: %v", err)
		return
	}

	_, err = stdin.Write([]byte(body))
	if err != nil {
		log.Printf("SendEmail3: %v", err)
		return
	}

	err = stdin.Close()
	if err != nil {
		log.Printf("SendEmail4: %v", err)
		return
	}

	err = cmd.Wait()
	if err != nil {
		log.Printf("SendEmail5: %v", err)
		return
	}

	return
}

func IsHidden(r *http.Request) bool {
	_, x := path.Split(r.URL.Path)
	return strings.Contains(x, "hidden")
}
