package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Identifier interface {
	Identity() int
	Name() string
}

type Comment struct {
	Id           int
	Stamp        time.Time
	Name         string
	Text         string
	Approved     bool
	ApprovalCode string
}

type Appearance struct {
	Subject         Identifier
	Feature         *Feature
	Computer        *Computer
	Description     Markup
	RealismStars    int
	Realism         Markup
	ImportanceStars int
	Importance      Markup
	VisibilityStars int
	Visibility      Markup
	Images          []string
	Comments        []Comment
}

type StarsInfo struct {
	LabelAlt   string
	LabelImage string
	StarsAlt   string
	StarsImage string
	Text       Markup
}

func MakeStarsInfo(label string, stars int, txt Markup) StarsInfo {
	var si StarsInfo
	si.Text = txt
	si.LabelAlt = label + ":"
	si.LabelImage = strings.ToLower(label) + ".png"
	si.StarsAlt = strings.Repeat("*", stars)
	if stars == 1 {
		si.StarsImage = "1star.png"
	} else {
		si.StarsImage = fmt.Sprintf("%dstars.png", stars)
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
		log.Print(err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var file string
		err = rows.Scan(&file)
		if err != nil {
			log.Print(err)
			continue
		}
		result = append(result, file)
	}

	return result
}

func (stc *Stc) AppearanceComments(computer, feature int) []Comment {
	result := []Comment{}

	rows, err := stc.Db.Query("SELECT "+
		"id, stamp, name, text, approved, approval_code"+
		" FROM comment WHERE approved=1 AND"+
		" feature=? AND computer=?",
		feature, computer)
	if err != nil {
		log.Print(err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var c Comment
		err = rows.Scan(&c.Id, &c.Stamp, &c.Name, &c.Text,
			&c.Approved, &c.ApprovalCode)
		if err != nil {
			log.Print(err)
			continue
		}
		result = append(result, c)
	}

	return result
}

func (stc *Stc) FeatureAppearances(f *Feature,
	hidden bool) ([]Appearance, error) {
	result := []Appearance{}

	rows, err := stc.Db.Query("SELECT "+
		"a.computer, a.description, "+
		"a.importance_stars, a.importance, "+
		"a.realism_stars, a.realism, "+
		"a.visibility_stars, a.visibility "+
		"FROM appearance AS a "+
		"LEFT JOIN computer AS c ON a.computer = c.id "+
		"WHERE feature=? AND visible=? "+
		"ORDER BY c.manufacturer, c.model",
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
		a.Comments = stc.AppearanceComments(a.Computer.Id, a.Feature.Id)
		result = append(result, a)
	}

	return result, nil
}

func (stc *Stc) ComputerAppearances(c *Computer,
	hidden bool) ([]Appearance, error) {
	result := []Appearance{}

	rows, err := stc.Db.Query("SELECT "+
		"a.feature, a.description, "+
		"a.importance_stars, a.importance, "+
		"a.realism_stars, a.realism, "+
		"a.visibility_stars, a.visibility "+
		"FROM appearance AS a "+
		"LEFT JOIN feature AS f ON a.feature = f.id "+
		"LEFT JOIN tv ON a.feature = tv.feature "+
		"WHERE computer=? AND visible=? "+
		"ORDER BY f.title, tv.season, tv.episode, tv.title",
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
		a.Comments = stc.AppearanceComments(a.Computer.Id, a.Feature.Id)
		result = append(result, a)
	}

	return result, nil
}

func (stc *Stc) LoadAppearance(fId, cId int) (*Appearance, error) {
	a := &Appearance{}

	err := stc.Db.QueryRow("SELECT  description, "+
		"importance_stars, importance, "+
		"realism_stars, realism, "+
		"visibility_stars, visibility "+
		"FROM appearance WHERE "+
		"feature=? AND computer=? AND visible=true",
		fId, cId).Scan(&a.Description,
		&a.ImportanceStars, &a.Importance,
		&a.RealismStars, &a.Realism,
		&a.VisibilityStars, &a.Visibility)
	if err != nil {
		return nil, err
	}

	a.Feature, err = stc.LoadFeature(fId)
	if err != nil {
		return nil, err
	}

	a.Computer, err = stc.LoadComputer(cId)
	if err != nil {
		return nil, err
	}
	a.Images = stc.AppearanceImages(a.Computer.Id, a.Feature.Id)
	a.Comments = stc.AppearanceComments(a.Computer.Id, a.Feature.Id)

	return a, nil
}

func (stc *Stc) AppearanceHandler(w http.ResponseWriter, r *http.Request) {
	form := SimpleForm(r)

	fId, err := strconv.Atoi(form["f"])
	if err != nil {
		http.Error(w, "bad feature id", 400)
		return
	}

	cId, err := strconv.Atoi(form["c"])
	if err != nil {
		http.Error(w, "bad computer id", 400)
		return
	}

	a, err := stc.LoadAppearance(fId, cId)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "bad appearance id", 400)
		return
	}

	title := PageTitle(a.Computer.Name() + " in " + a.Feature.Name())

	err = stc.Template.Exec("appearance", w, AppearanceTemplateData{
		PageTitle:  title,
		Appearance: a,
		Form:       stc.ParseCommentForm(form, a),
	})
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "bad appearance", 500)
		return
	}

}
