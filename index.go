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

type IndexPoint struct {
	Name string
	Link string
}

type Index struct {
	Indices []IndexPoint
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

	stc.ComputersByManufacturer, err = stc.LoadComputersByManufacturer()
	if err != nil {
		return err
	}

	return nil
}

func (i *Index) addIndexPoint(name, link string) {
	curLen := len(i.Indices)
	if curLen == 0 || i.Indices[curLen-1].Name != name {
		i.Indices = append(i.Indices, IndexPoint{
			Name: name,
			Link: link,
		})
	}
}

func (i *Index) addAlphaIndexPoint(name string) {
	curLen := len(i.Indices)
	for ix, n := range "0ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		in := string(n)
		if curLen > ix {
			if in == name {
				break
			}
			continue
		}
		if in == name {
			i.addIndexPoint(name, name)
			break
		}
		i.addIndexPoint(in, "")
	}
}

func (stc *Stc) LoadFeaturesByName() (*Index, error) {
	log.Print("Loading FeaturesByName index")

	i := &Index{
		Indices: []IndexPoint{},
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
		if len(appears) == 0 {
			continue
		}
		things := ""
		for _, a := range appears {
			things += NonBroken("• "+a.Computer.Name()) + " "
		}

		index := IndexChar(f.Title)
		i.addAlphaIndexPoint(index)

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
		Indices: []IndexPoint{},
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
		if len(appears) == 0 {
			continue
		}
		things := ""
		for _, a := range appears {
			things += NonBroken("• "+a.Computer.Name()) + " "
		}

		index := fmt.Sprintf("%d", (f.Year/10)*10)
		i.addIndexPoint(index, index)

		i.Entries[index] = append(i.Entries[index], IndexItem{
			Name:   f.Name(),
			Link:   fmt.Sprintf("/feature.html?f=%d", f.Id),
			Things: things,
		})
	}

	return i, nil
}

func (stc *Stc) LoadComputersByManufacturer() (*Index, error) {
	log.Print("Loading ComputersByManufacturer index")

	i := &Index{
		Indices: []IndexPoint{},
		AltName: "",
		AltLink: "",
		Entries: map[string][]IndexItem{},
	}

	rows, err := stc.Db.Query("SELECT id FROM computer " +
		"ORDER BY manufacturer, model")
	if err != nil {
		log.Printf("LoadComputersByManufacturer1: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var compId int
		err = rows.Scan(&compId)
		if err != nil {
			log.Printf("LoadComputersByManufacturer2: %v", err)
			return nil, err
		}

		c, err := stc.LoadComputer(compId)
		if err != nil {
			log.Printf("LoadComputersByManufacturer3 (%d): %v", compId, err)
			return nil, err
		}

		appears, err := stc.ComputerAppearances(c, false)
		if err != nil {
			log.Printf("LoadComputersByManufacturer4: %v", err)
			return nil, err
		}
		if len(appears) == 0 {
			continue
		}
		things := ""
		for _, a := range appears {
			things += NonBroken("• "+a.Feature.Name()) + " "
		}

		index := IndexChar(c.Manufacturer)
		i.addAlphaIndexPoint(index)

		i.Entries[index] = append(i.Entries[index], IndexItem{
			Name:   c.Name(),
			Link:   fmt.Sprintf("/computer.html?c=%d", c.Id),
			Things: things,
		})
	}

	return i, nil
}

func (stc *Stc) IndexHandler(index *Index) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := stc.Template.Exec("index", w, &IndexTemplate{
			PageTitle: PageTitle(""),
			Index:     index,
		})
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "bad index", 500)
			return
		}
	}
}
