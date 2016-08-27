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

	stc.FeaturesByYear, err = stc.LoadFeaturesByYear()
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

	rows, err := stc.Db.Query("SELECT feature.id FROM " +
		"feature LEFT JOIN tv ON feature.id = tv.feature " +
		"ORDER BY feature.title,tv.season,tv.episode,tv.title")
	if err != nil {
		log.Printf("LoadFeaturesByName1: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var featId int
		err = rows.Scan(&featId)
		if err != nil {
			log.Printf("LoadFeaturesByName2: %v", err)
			return nil, err
		}

		f, err := stc.LoadFeature(featId)
		if err != nil {
			log.Printf("LoadFeaturesByName3 (%d): %v", featId, err)
			return nil, err
		}

		appears, err := stc.FeatureAppearances(f, false)
		if err != nil {
			log.Printf("LoadFeaturesByName4: %v", err)
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

func (stc *Stc) LoadFeaturesByYear() (*Index, error) {
	log.Print("Loading FeaturesByYear index")

	i := &Index{
		Indices: []string{},
		AltName: "name",
		AltLink: "features.html",
		Entries: map[string][]IndexItem{},
	}

	rows, err := stc.Db.Query("SELECT feature.id FROM " +
		"feature LEFT JOIN tv ON feature.id = tv.feature " +
		"ORDER BY feature.year,feature.title,tv.season,tv.episode,tv.title")
	if err != nil {
		log.Printf("LoadFeaturesByYear1: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var featId int
		err = rows.Scan(&featId)
		if err != nil {
			log.Printf("LoadFeaturesByYear2: %v", err)
			return nil, err
		}

		f, err := stc.LoadFeature(featId)
		if err != nil {
			log.Printf("LoadFeaturesByYear3 (%d): %v", featId, err)
			return nil, err
		}

		appears, err := stc.FeatureAppearances(f, false)
		if err != nil {
			log.Printf("LoadFeaturesByYear4: %v", err)
			return nil, err
		}
		things := ""
		for _, a := range appears {
			things += NonBroken("• "+a.Computer.Name()) + " "
		}

		index := fmt.Sprintf("%d", (f.Year/10)*10)
		if len(i.Indices) == 0 ||
			i.Indices[len(i.Indices)-1] != index {
			i.Indices = append(i.Indices, index)
		}

		i.Entries[index] = append(i.Entries[index], IndexItem{
			Name:   f.Name(),
			Link:   fmt.Sprintf("/feature.html?f=%d", f.Id),
			Things: things,
		})
	}

	return i, nil
}

func (stc *Stc) MakeIndexHandler(index *Index) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := stc.Template.Exec("index", w, &IndexTemplate{
			PageTitle: "Starring the Computer",
			Index:     index,
		})
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "bad index", 500)
			return
		}
	}
}
