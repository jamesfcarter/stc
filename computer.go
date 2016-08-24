package main

import "fmt"

type Computer struct {
	Stc          *Stc
	Id           int
	Manufacturer string
	Model        string
	Description  string
	Image        string
	InfoLink     string
}

func (stc *Stc) LoadComputer(id int) (*Computer, error) {
	c := &Computer{}

	err := stc.Db.QueryRow("SELECT "+
		"manufacturer, model, description, image, info_link"+
		" FROM computer WHERE id=?", id).Scan(&c.Manufacturer,
		&c.Model, &c.Description, &c.Image, &c.InfoLink)
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
		Id:          c.Id,
		Image:       c.Image,
		Name:        c.Name(),
		InfoLink:    c.InfoLink,
		Description: c.Description,
		Appearances: appearances,
	}
}

func (c *Computer) Identity() int {
	return c.Id
}

func (c *Computer) Name() string {
	return fmt.Sprintf("%s %s", c.Manufacturer, c.Model)
}
