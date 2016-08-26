package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

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
		Id:          f.Id,
		Title:       f.Title,
		Image:       fmt.Sprintf("%d.jpg", f.Id),
		Name:        f.Title,
		ImdbLink:    f.ImdbLink,
		Description: f.Description,
		Appearances: appearances,
	}
}

func (f *Feature) Identity() int {
	return f.Id
}

func (f *Feature) Name() string {
	return f.Title
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
