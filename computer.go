package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Computer struct {
	Stc          *Stc
	Id           int
	Manufacturer string
	Model        string
	Description  Markup
	Image        string
	InfoLink     string
	ImageLink    sql.NullString
}

func (stc *Stc) LoadComputer(id int) (*Computer, error) {
	c := &Computer{}

	err := stc.Db.QueryRow("SELECT "+
		"manufacturer, model, description, image, info_link, image_link"+
		" FROM computer WHERE id=?", id).Scan(&c.Manufacturer,
		&c.Model, &c.Description, &c.Image, &c.InfoLink, &c.ImageLink)
	if err != nil {
		return nil, err
	}
	c.Id = id
	c.Stc = stc
	return c, nil
}

func (c *Computer) TemplateData(deep, hidden bool) ComputerTemplateData {
	var appearances []Appearance
	if deep {
		appearances, _ = c.Stc.ComputerAppearances(c, hidden)
	}
	return ComputerTemplateData{
		PageTitle:   PageTitle(c.Name()),
		Computer:    c,
		Appearances: appearances,
	}
}

func (c *Computer) Identity() int {
	return c.Id
}

func (c *Computer) Name() string {
	return fmt.Sprintf("%s %s", c.Manufacturer, c.Model)
}

func (stc *Stc) ComputerHandler(w http.ResponseWriter, r *http.Request) {
	form := SimpleForm(r)
	hidden := IsHidden(r)

	id, err := strconv.Atoi(form["c"])
	if err != nil {
		http.Error(w, "bad computer id", 400)
		return
	}
	c, err := stc.LoadComputer(id)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "bad computer id", 400)
		return
	}
	err = stc.Template.Exec("computer", w, c.TemplateData(true, hidden))
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "bad feature", 500)
		return
	}
}
