package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
)

type Stc struct {
	Db       *sql.DB
	Template *Templates
}

func main() {
	stc := &Stc{}

	// user:pass@tcp(host:3306)/database
	dbSpec := os.Getenv("STC_DB")
	if dbSpec == "" {
		log.Fatal("STC_DB not set")
	}
	var err error
	stc.Db, err = sql.Open("mysql", dbSpec)
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

	http.Handle("/movies/", fs)
	http.Handle("/computers/", fs)
	http.Handle("/snapshots/", fs)
	http.Handle("/unprocessed/", fs)
	http.Handle("/img/", fs)
	http.Handle("/favicon.ico", fs)

	http.HandleFunc("/feature.html", stc.FeatureHandler)
	http.HandleFunc("/stylesheet.css", func(w http.ResponseWriter, r *http.Request) {
		err = stc.Template.Exec("stylesheet", w, nil)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "bad stylesheet", 500)
			return
		}
	})
	http.ListenAndServe(":8080", nil)
}
