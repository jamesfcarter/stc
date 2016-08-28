package main

import (
	"html/template"
	"strconv"
	"strings"
)

type CommentForm struct {
	Name    string
	Comment string
	Year    string
	Errors  map[string]string
	Created bool
}

func (cf *CommentForm) Label(label string) template.HTML {
	tag := strings.ToLower(strings.Split(label, " ")[0])
	if cf.Errors[tag] == "" {
		return template.HTML(label)
	}
	label = strings.TrimSpace(strings.Split(label, "(")[0])
	errmsg := " <span class=\"error\">" + cf.Errors[tag] + "</span>"
	return template.HTML(label + errmsg + ":")
}

func (cf *CommentForm) Valid() bool {
	return len(cf.Errors) == 0
}

func (cf *CommentForm) Empty() bool {
	return cf.Name == "" &&
		cf.Comment == "" &&
		cf.Year == ""
}

func ParseCommentForm(form map[string]string, feat *Feature) *CommentForm {
	cf := &CommentForm{
		Name:    form["n"],
		Comment: form["t"],
		Year:    form["y"],
		Errors:  map[string]string{},
	}
	if cf.Empty() {
		return cf
	}

	if cf.Name == "" {
		cf.Errors["name"] = "missing!"
	}
	if cf.Comment == "" {
		cf.Errors["comment"] = "empty!"
	}
	if cf.Year == "" {
		cf.Errors["year"] = "missing!"
	} else {
		year, err := strconv.Atoi(cf.Year)
		if err != nil || year != feat.Year {
			cf.Errors["year"] = "does not match!"
		}
	}

	if cf.Valid() {
		// FIXME: post comment
		return &CommentForm{
			Errors:  map[string]string{},
			Created: true,
		}
	}

	return cf
}
