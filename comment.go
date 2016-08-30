package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
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

	_, err := a.Feature.Stc.Db.Exec("INSERT comment SET "+
		"feature=?, computer=?, name=?, text=?, "+
		"approved=?, approval_code=?",
		a.Feature.Id, a.Computer.Id,
		cf.Name, cf.Comment, false, approvalCode)
	if err != nil {
		log.Printf("unable to store comment: %v")
		return
	}

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

func (stc *Stc) CommentURL(code int) string {
	var fId, cId int

	_ = stc.Db.QueryRow("SELECT feature, computer FROM "+
		"comment WHERE approval_code=?", code).Scan(&fId, &cId)
	return fmt.Sprintf("/appearance.html?f=%d&c=%d", fId, cId)
}

func (stc *Stc) CommentApprove(code int) error {
	res, err := stc.Db.Exec("UPDATE comment SET approved=? WHERE "+
		"approval_code=?", true, code)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		return fmt.Errorf("Could not find comment to approve.")
	}
	return nil
}

func (stc *Stc) CommentDelete(code int) error {
	res, err := stc.Db.Exec("DELETE FROM comment WHERE "+
		"approval_code=?", code)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		return fmt.Errorf("Could not find comment to delete.")
	}
	return nil
}

func (stc *Stc) CommentHandler(w http.ResponseWriter, r *http.Request) {
	action, codeS := path.Split(r.URL.Path)

	code, err := strconv.Atoi(codeS)
	if err != nil {
		http.Error(w, "bad comment code", 400)
		return
	}

	// Prevent a an attacker cycling through all the approval codes
	// in a reasonable time.
	_ = <-stc.ApprovalQueue
	go func() {
		time.Sleep(time.Second)
		stc.ApprovalQueue <- struct{}{}
	}()

	switch action {
	case "/approve/":
		err = stc.CommentApprove(code)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), 400)
			return
		}
		err = stc.Template.Exec("approvecomment", w, struct {
			PageTitle string
			Link      string
		}{
			PageTitle: PageTitle("Approved Comment"),
			Link:      stc.CommentURL(code),
		})
	case "/deny/":
		err = stc.Template.Exec("denycomment", w, struct {
			PageTitle string
			Link      string
			DelLink   string
		}{
			PageTitle: PageTitle("Deny Comment"),
			Link:      stc.CommentURL(code),
			DelLink:   fmt.Sprintf("/delete/%d", code),
		})
	case "/delete/":
		err = stc.CommentDelete(code)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), 400)
			return
		}
		err = stc.Template.Exec("deletecomment", w, struct {
			PageTitle string
		}{
			PageTitle: PageTitle("Comment Deleted"),
		})
	}
}
