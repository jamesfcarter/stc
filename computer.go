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
		" FROM feature WHERE id=?", id).Scan(&c.Manufacturer,
		&c.Model, &c.Description, &c.Image, &c.InfoLink)
	if err != nil {
		return nil, err
	}
	c.Id = id
	c.Stc = stc
	return c, nil
}

func (c *Computer) TemplateData(deep bool) ComputerTemplateData {
	var finfo []FeatureTemplateData
	if deep {
		features, _ := c.Stc.FeaturesForComputer(c.Id)
		finfo = make([]FeatureTemplateData, len(features))
		for i, v := range features {
			finfo[i] = v.TemplateData(false)
		}
	}
	return ComputerTemplateData{
		Id:          c.Id,
		Image:       c.Image,
		Name:        fmt.Sprintf("%s %s", c.Manufacturer, c.Model),
		InfoLink:    c.InfoLink,
		Description: c.Description,
		Features:    finfo,
	}
}
