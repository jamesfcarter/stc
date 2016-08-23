package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
)

var Db *sql.DB

func main() {
	var err error
	dbSpec := os.Getenv("STC_DB")
	if dbSpec == "" {
		log.Fatal("STC_DB not set")
	}
	Db, err = sql.Open("mysql", dbSpec)
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()
	err = Db.Ping()
	if err != nil {
		log.Fatal("can't communicate with db: %v", err)
	}

	fsRoot := os.Getenv("STC_ROOT")
	if fsRoot == "" {
		panic("STC_ROOT not set")
	}
	fs := http.FileServer(http.Dir(fsRoot))

	http.Handle("/movies/", fs)
	http.Handle("/computers/", fs)
	http.Handle("/snapshots/", fs)
	http.Handle("/unprocessed/", fs)

	http.HandleFunc("/feature.html", FeatureHandler)
	http.ListenAndServe(":8080", nil)
}
