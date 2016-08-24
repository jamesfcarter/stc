package main

import "database/sql"

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
