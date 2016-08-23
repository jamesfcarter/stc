package main

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
