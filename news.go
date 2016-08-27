package main

import (
	"fmt"
	"log"
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
