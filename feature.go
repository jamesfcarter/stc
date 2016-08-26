package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Tv struct {
	Season  sql.NullInt64
	Episode sql.NullInt64
	Title   sql.NullString
}

func (stc *Stc) LoadTv(id int) (*Tv, error) {
	t := &Tv{}

	err := stc.Db.QueryRow("SELECT "+
		"season, episode, title"+
		" FROM tv WHERE feature=?", id).Scan(&t.Season, &t.Episode, &t.Title)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Tv) FullName() string {
	parts := []string{}
	if t.Season.Valid {
		parts = append(parts,
			fmt.Sprintf("Season %d", t.Season.Int64))
	}
	if t.Episode.Valid {
		parts = append(parts,
			fmt.Sprintf("Episode %d", t.Episode.Int64))
	}
	if t.Title.Valid {
		parts = append(parts,
			"\""+ReadableTitle(t.Title.String)+"\"")
	}
	return strings.Join(parts, ", ")
}

type Feature struct {
	Stc         *Stc
	Id          int
	Title       string
	IsTvEpisode bool
	Year        int
	ImdbLink    string
	Description string
}

func (stc *Stc) LoadFeature(id int) (*Feature, error) {
	f := &Feature{}

	err := stc.Db.QueryRow("SELECT "+
		"title, is_tv_episode, year, imdb_link, description"+
		" FROM feature WHERE id=?", id).Scan(&f.Title, &f.IsTvEpisode,
		&f.Year, &f.ImdbLink, &f.Description)
	if err != nil {
		return nil, err
	}
	f.Id = id
	f.Stc = stc
	return f, nil
}

func (f *Feature) TemplateData(deep, hidden bool) FeatureTemplateData {
	var appearances []Appearance
	var err error
	if deep {
		appearances, err = f.Stc.FeatureAppearances(f, hidden)
		if err != nil {
			log.Printf("%v", err)
		}
	}
	return FeatureTemplateData{
		PageTitle:   PageTitle(ReadableTitle(f.Title)),
		Feature:     f,
		Appearances: appearances,
	}
}

func (f *Feature) Identity() int {
	return f.Id
}

func (f *Feature) Image() string {
	return fmt.Sprintf("%d.jpg", f.Id)
}

func (f *Feature) Name() string {
	name := ReadableTitle(f.Title)
	if f.IsTvEpisode {
		tv, err := f.Stc.LoadTv(f.Id)
		if err != nil {
			log.Print(err)
		} else {
			name += " - " + tv.FullName()
		}
	}
	return name
}

func (stc *Stc) FeatureHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("f"))
	if err != nil {
		http.Error(w, "bad feature id", 400)
		return
	}
	f, err := stc.LoadFeature(id)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "bad feature id", 400)
		return
	}
	err = stc.Template.Exec("feature", w, f.TemplateData(true, false))
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "bad feature", 500)
		return
	}
}
