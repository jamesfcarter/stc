package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Stc struct {
	Db                      *sql.DB
	Template                *Templates
	Root                    string
	Film                    *Film
	FeaturesByName          *Index
	FeaturesByYear          *Index
	ComputersByManufacturer *Index
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	stc := &Stc{}

	endpoint := os.Getenv("STC_ENDPOINT")
	if endpoint == "" {
		endpoint = ":8080"
	}

	// user:pass@tcp(host:3306)/database
	dbSpec := os.Getenv("STC_DB")
	if dbSpec == "" {
		log.Fatal("STC_DB not set")
	}
	var err error
	stc.Db, err = sql.Open("mysql", dbSpec+"?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer stc.Db.Close()
	err = stc.Db.Ping()
	if err != nil {
		log.Fatalf("can't communicate with db: %v", err)
	}

	stc.Root = os.Getenv("STC_ROOT")
	if stc.Root == "" {
		panic("STC_ROOT not set")
	}
	fs := http.FileServer(http.Dir(stc.Root))

	stc.Template, err = MakeTemplates()
	if err != nil {
		log.Fatalf("could not make templates: %v", err)
	}

	err = stc.LoadIndices()
	if err != nil {
		log.Fatalf("could not load indices: %v", err)
	}

	stc.Film, err = stc.NewFilm()
	if err != nil {
		log.Fatalf("could not build film image: %v", err)
	}

	http.Handle("/movies/", fs)
	http.Handle("/computers/", fs)
	http.Handle("/snapshots/", fs)
	http.Handle("/unprocessed/", fs)
	http.Handle("/img/", fs)
	http.Handle("/favicon.ico", fs)

	http.HandleFunc("/film.jpg", stc.FilmHandler)

	http.HandleFunc("/feature.html", stc.FeatureHandler)
	http.HandleFunc("/computer.html", stc.ComputerHandler)

	http.HandleFunc("/features.html",
		stc.MakeIndexHandler(stc.FeaturesByName))
	http.HandleFunc("/featuresyear.html",
		stc.MakeIndexHandler(stc.FeaturesByYear))
	http.HandleFunc("/computers.html",
		stc.MakeIndexHandler(stc.ComputersByManufacturer))

	http.HandleFunc("/stylesheet.css", func(w http.ResponseWriter, r *http.Request) {
		err = stc.Template.Exec("stylesheet", w, nil)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "bad stylesheet", 500)
			return
		}
	})

	log.Printf("Starting service on %s", endpoint)
	http.ListenAndServe(endpoint, nil)
}
