package main

import (
	"fmt"
	"html/template"
	"math/rand"
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

func (stc *Stc) ParseCommentForm(form map[string]string, a *Appearance) *CommentForm {
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
		if err != nil || year != a.Feature.Year {
			cf.Errors["year"] = "does not match!"
		}
	}

	if cf.Valid() {
		cf.Post(a)
		return &CommentForm{
			Errors:  map[string]string{},
			Created: true,
		}
	}

	return cf
}

func (cf *CommentForm) Post(a *Appearance) {
	approvalCode := rand.Intn(1000000000)
	stc := "http://starringthecomputer.com/"
	linkApprove := fmt.Sprintf("%sapprove/%d", stc, approvalCode)
	linkDeny := fmt.Sprintf("%sdeny/%d", stc, approvalCode)
	link := fmt.Sprintf("%sappearance.html?f=%d&c=%d",
		stc, a.Feature.Id, a.Computer.Id)

	// FIXME: create comment

	subject := "STC Comment from " + cf.Name
	msg := "Name: " + cf.Name + "\n" +
		"Comment: " + cf.Comment + "\n\n" +
		"Feature: " + a.Feature.Name() + "\n" +
		"Computer: " + a.Computer.Name() + "\n\n" +
		link + "\n\n" +
		"APPROVE: " + linkApprove + "\n" +
		"DENY: " + linkDeny + "\n\n"

	SendEmail(subject, msg)
}
