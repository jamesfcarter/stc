package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type News struct {
	Stc   *Stc
	Title string
	Text  Markup
	Stamp time.Time
}

func (stc *Stc) LoadNews(offset, count int) ([]News, error) {
	r := []News{}

	rows, err := stc.Db.Query("SELECT title, text, stamp FROM news " +
		"ORDER BY stamp DESC LIMIT " +
		fmt.Sprintf("%d, %d", offset, count))
	if err != nil {
		log.Printf("LoadNews1: %v", err)
		return r, err
	}
	defer rows.Close()

	for rows.Next() {
		var n News

		err = rows.Scan(&n.Title, &n.Text, &n.Stamp)
		if err != nil {
			log.Printf("LoadNews1: %v", err)
			return r, err
		}
		r = append(r, n)
	}

	return r, nil
}

func (stc *Stc) NewsTotal() (int, error) {
	var total int

	err := stc.Db.QueryRow("SELECT COUNT(*) FROM news").Scan(&total)
	return total, err
}

func (stc *Stc) NewsHandler(w http.ResponseWriter, r *http.Request) {
	offset, err := strconv.Atoi(r.URL.Query().Get("o"))
	if err != nil {
		offset = 0
	}
	count, err := strconv.Atoi(r.URL.Query().Get("s"))
	if err != nil {
		count = 16
	}
	news, err := stc.LoadNews(offset, count)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "server error", 500)
		return
	}
	total, err := stc.NewsTotal()
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "server error", 500)
		return
	}
	newer := ""
	if offset > 0 {
		newo := offset - count
		if newo >= 0 {
			newer = fmt.Sprintf("o=%d&s=%d", newo, count)
		}
	}
	older := ""
	if offset <= total-count {
		oldo := offset + count
		if oldo <= total-count+1 {
			older = fmt.Sprintf("o=%d&s=%d", oldo, count)
		}
	}
	err = stc.Template.Exec("news", w, NewsTemplateData{
		PageTitle: "Starring the Computer",
		LinkNewer: newer,
		LinkOlder: older,
		News:      news,
	})
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "server error", 500)
		return
	}
}
