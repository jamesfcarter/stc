package main

import (
	"net/http"
	"os"
)

func main() {
	fsRoot := os.Getenv("STC_ROOT")
	if fsRoot == "" {
		panic("STC_ROOT not set")
	}
	fs := http.FileServer(http.Dir(fsRoot))

	http.Handle("/movies/", fs)
	http.Handle("/computers/", fs)
	http.Handle("/snapshots/", fs)
	http.Handle("/unprocessed/", fs)
	http.ListenAndServe(":8080", nil)
}
