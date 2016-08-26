package main

import "database/sql"
import "strings"

type Identifier interface {
	Identity() int
	Name() string
}

type Appearance struct {
	Subject         Identifier
	Feature         *Feature
	Computer        *Computer
	Description     string
	RealismStars    int
	Realism         sql.NullString
	ImportanceStars int
	Importance      sql.NullString
	VisibilityStars int
	Visibility      sql.NullString
	Images          []string
}

type StarsInfo struct {
	LabelAlt   string
	LabelImage string
	StarsAlt   string
	StarsImage string
	Text       string
}

func MakeStarsInfo(label string, stars int, txt sql.NullString) StarsInfo {
	var si StarsInfo
	si.Text = txt.String
	si.LabelAlt = label + ":"
	si.LabelImage = strings.ToLower(label) + ".png"
	switch {
	case stars == 1:
		si.StarsAlt = "*"
		si.StarsImage = "1star.png"
	case stars == 2:
		si.StarsAlt = "**"
		si.StarsImage = "2stars.png"
	case stars == 3:
		si.StarsAlt = "***"
		si.StarsImage = "3stars.png"
	case stars == 4:
		si.StarsAlt = "****"
		si.StarsImage = "4stars.png"
	case stars == 5:
		si.StarsAlt = "*****"
		si.StarsImage = "5stars.png"
	}
	return si
}

func (a Appearance) RealismInfo() StarsInfo {
	return MakeStarsInfo("Realism", a.RealismStars, a.Realism)
}

func (a Appearance) ImportanceInfo() StarsInfo {
	return MakeStarsInfo("Importance", a.ImportanceStars, a.Importance)
}

func (a Appearance) VisibilityInfo() StarsInfo {
	return MakeStarsInfo("Visibility", a.VisibilityStars, a.Visibility)
}

func (stc *Stc) AppearanceImages(computer, feature int) []string {
	result := []string{}

	rows, err := stc.Db.Query("SELECT "+
		"file FROM image WHERE feature=? AND computer=?",
		feature, computer)
	if err != nil {
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var file string
		err = rows.Scan(&file)
		if err != nil {
			continue
		}
		result = append(result, file)
	}

	return result
}

func (stc *Stc) FeatureAppearances(f *Feature,
	hidden bool) ([]Appearance, error) {
	result := []Appearance{}

	rows, err := stc.Db.Query("SELECT "+
		"computer, description, "+
		"importance_stars, importance, "+
		"realism_stars, realism, "+
		"visibility_stars, visibility "+
		"FROM appearance WHERE feature=? AND visible=?",
		f.Id, !hidden)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var a Appearance
		err = rows.Scan(&cid, &a.Description,
			&a.ImportanceStars, &a.Importance,
			&a.RealismStars, &a.Realism,
			&a.VisibilityStars, &a.Visibility)
		if err != nil {
			return nil, err
		}
		a.Feature = f
		a.Computer, err = stc.LoadComputer(cid)
		if err != nil {
			return nil, err
		}
		a.Subject = a.Computer
		a.Images = stc.AppearanceImages(a.Computer.Id, a.Feature.Id)
		result = append(result, a)
	}

	return result, nil
}

func (stc *Stc) ComputerAppearances(c *Computer,
	hidden bool) ([]Appearance, error) {
	result := []Appearance{}

	rows, err := stc.Db.Query("SELECT "+
		"feature, description, "+
		"importance_stars, importance, "+
		"realism_stars, realism, "+
		"visibility_stars, visibility "+
		"FROM appearance WHERE computer=? AND visible=?",
		c.Id, !hidden)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var fid int
		var a Appearance
		err = rows.Scan(&fid, &a.Description,
			&a.ImportanceStars, &a.Importance,
			&a.RealismStars, &a.Realism,
			&a.VisibilityStars, &a.Visibility)
		if err != nil {
			return nil, err
		}
		a.Computer = c
		a.Feature, err = stc.LoadFeature(fid)
		if err != nil {
			return nil, err
		}
		a.Subject = a.Feature
		a.Images = stc.AppearanceImages(a.Computer.Id, a.Feature.Id)
		result = append(result, a)
	}

	return result, nil
}
