package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
)

type Stc struct {
	Db             *sql.DB
	Template       *Templates
	FeaturesByName *Index
}

func main() {
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

	fsRoot := os.Getenv("STC_ROOT")
	if fsRoot == "" {
		panic("STC_ROOT not set")
	}
	fs := http.FileServer(http.Dir(fsRoot))

	stc.Template, err = MakeTemplates()
	if err != nil {
		log.Fatalf("could not make templates: %v", err)
	}

	err = stc.LoadIndices()
	if err != nil {
		log.Fatalf("could not load indices: %v", err)
	}

	http.Handle("/movies/", fs)
	http.Handle("/computers/", fs)
	http.Handle("/snapshots/", fs)
	http.Handle("/unprocessed/", fs)
	http.Handle("/img/", fs)
	http.Handle("/favicon.ico", fs)

	http.HandleFunc("/feature.html", stc.FeatureHandler)
	http.HandleFunc("/computer.html", stc.ComputerHandler)
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
