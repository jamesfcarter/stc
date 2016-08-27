package main

import (
	"fmt"
	"log"
	"net/http"
)

type IndexItem struct {
	Name   string
	Link   string
	Things string
}

type Index struct {
	Indices []string
	AltName string
	AltLink string
	Entries map[string][]IndexItem
}

func AToZ() []string {
	return []string{
		"0", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K",
		"L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W",
		"X", "Y", "Z",
	}
}

func (stc *Stc) LoadIndices() error {
	var err error

	stc.FeaturesByName, err = stc.LoadFeaturesByName()
	if err != nil {
		return err
	}

	return nil
}

func (stc *Stc) LoadFeaturesByName() (*Index, error) {
	log.Print("Loading FeaturesByName index")

	i := &Index{
		Indices: AToZ(),
		AltName: "year",
		AltLink: "featuresyear.html",
		Entries: map[string][]IndexItem{},
	}

	rows, err := stc.Db.Query("SELECT feature.id FROM feature,tv " +
		"WHERE feature.id = tv.feature ORDER BY feature.title," +
		"tv.season,tv.episode,tv.title")
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var featId int
		err = rows.Scan(&featId)
		if err != nil {
			log.Print(err)
			return nil, err
		}

		f, err := stc.LoadFeature(featId)
		if err != nil {
			log.Print(err)
			return nil, err
		}

		appears, err := stc.FeatureAppearances(f, false)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		things := ""
		for _, a := range appears {
			things += NonBroken("• "+a.Computer.Name()) + " "
		}

		index := IndexChar(f.Title)

		i.Entries[index] = append(i.Entries[index], IndexItem{
			Name:   f.Name(),
			Link:   fmt.Sprintf("/feature.html?f=%d", f.Id),
			Things: things,
		})
	}

	return i, nil
}

func (stc *Stc) FeaturesHandler(w http.ResponseWriter, r *http.Request) {
	err := stc.Template.Exec("index", w, &IndexTemplate{
		PageTitle: "",
		Index:     stc.FeaturesByName,
	})
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "bad index", 500)
		return
	}
}
