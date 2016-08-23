package main

import (
	"log"
	"net/http"
	"strconv"
)

type Feature struct {
	Id          int
	Title       string
	IsTvEpisode bool
	Year        int
	ImdbLink    string
	Description string
}

func LoadFeature(id int) (*Feature, error) {
	f := &Feature{}

	err := Db.QueryRow("SELECT "+
		"title, is_tv_episode, year, imdb_link, description"+
		" FROM feature WHERE id=?", id).Scan(&f.Title, &f.IsTvEpisode,
		&f.Year, &f.ImdbLink, &f.Description)
	if err != nil {
		return nil, err
	}
	f.Id = id
	return f, nil
}

func FeatureHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("f"))
	if err != nil {
		http.Error(w, "bad feature id", 400)
		return
	}
	f, err := LoadFeature(id)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "bad feature id", 400)
		return
	}
	log.Printf("%v", f)
}
