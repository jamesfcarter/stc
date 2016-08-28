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

func (stc *Stc) LoadNewsItem(item string) ([]News, error) {
	n := News{}

	err := stc.Db.QueryRow("SELECT title, text, stamp FROM news "+
		"WHERE stamp = ?", item).Scan(&n.Title, &n.Text, &n.Stamp)
	if err != nil {
		log.Printf("LoadNewsItem1: %v", err)
		return []News{}, err
	}

	return []News{n}, nil
}

func (stc *Stc) NewsTotal() (int, error) {
	var total int

	err := stc.Db.QueryRow("SELECT COUNT(*) FROM news").Scan(&total)
	return total, err
}

func (stc *Stc) NewsItemHandler(w http.ResponseWriter, r *http.Request) {
	form := SimpleForm(r)

	item := form["i"]
	if item == "" {
		http.Error(w, "bad item id", 400)
		return
	}

	news, err := stc.LoadNewsItem(item + ":00")
	if err != nil {
		http.Error(w, "item not found", 404)
		return
	}

	err = stc.Template.Exec("newsitem", w, NewsTemplateData{
		PageTitle: PageTitle(news[0].Title),
		LinkNewer: "",
		LinkOlder: "",
		News:      news,
	})
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "server error", 500)
		return
	}
}

func (stc *Stc) NewsHandler(w http.ResponseWriter, r *http.Request) {
	form := SimpleForm(r)

	offset, err := strconv.Atoi(form["o"])
	if err != nil {
		offset = 0
	}
	count, err := strconv.Atoi(form["s"])
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
		PageTitle: PageTitle(""),
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

func (stc *Stc) RssHandler(w http.ResponseWriter, r *http.Request) {
	news, err := stc.LoadNews(0, 10)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "server error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/rss+xml")
	w.Write([]byte("<?xml version='1.0' encoding='UTF-8'?>\n"))
	err = stc.Template.Exec("rss", w, RssTemplateData{
		Now:       time.Now(),
		RssFormat: "Mon, 02 Jan 2006 15:04:05 -0700",
		IndexTime: "2006-01-02 15:04",
		News:      news,
	})
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "server error", 500)
		return
	}
}
