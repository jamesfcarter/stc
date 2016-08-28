package main

import (
	"bytes"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	FilmWidth      = uint(150)
	FilmFrames     = 16
	FilmTimeout    = 60
	FilmBackground = "/img/film_bg.png"
)

type Film struct {
	Stc        *Stc
	Image      image.Image
	Made       time.Time
	Background image.Image
}

func (stc *Stc) NewFilm() (*Film, error) {
	var err error
	f := &Film{
		Stc: stc,
	}

	f.Background, err = f.loadPngImage(stc.Root+FilmBackground, FilmWidth)
	if err != nil {
		return nil, err
	}

	err = f.Update()
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (f *Film) loadPngImage(name string, width uint) (image.Image, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return resize.Resize(width, 0, img, resize.Lanczos3), nil
}

func (f *Film) loadJpegImage(name string, width uint) (image.Image, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}

	return resize.Resize(width, 0, img, resize.Lanczos3), nil
}

func (f *Film) randomImage(width uint) (image.Image, error) {
	files, err := ioutil.ReadDir(f.Stc.Root + "/snapshots")
	if err != nil {
		return nil, err
	}
	file := f.Stc.Root + "/snapshots/" +
		files[rand.Intn(len(files))].Name()
	return f.loadJpegImage(file, width)
}

func (f *Film) makeBackground(height uint) (draw.Image, uint) {
	filmHeight := uint(f.Background.Bounds().Dy())
	filmReps := int((height / filmHeight) + 1)

	r := image.NewRGBA(image.Rect(0, 0,
		int(FilmWidth), filmReps*int(filmHeight)))

	for i := 0; i < filmReps; i++ {
		startY := i * int(filmHeight)
		targetRect := image.Rect(0, startY,
			int(FilmWidth), startY+int(filmHeight))
		draw.Draw(r, targetRect, f.Background,
			image.Pt(0, 0), draw.Over)
	}

	extraGap := (uint(filmReps)*filmHeight - height) / FilmFrames

	return r, extraGap
}

func (f *Film) Update() error {
	var err error
	var height, imageWidth, gap uint

	log.Printf("Updating film image")

	imageWidth = FilmWidth * 7 / 10
	gap = uint(FilmWidth / 25)

	images := make([]image.Image, FilmFrames)
	for i := range images {
		// Some images fail to load so allow a few retries
		for n := 0; n < 10; n++ {
			images[i], err = f.randomImage(imageWidth)
			if err == nil {
				break
			}
		}
		if err != nil {
			return err
		}
		height += uint(images[i].Bounds().Dy()) + gap
	}

	film, extraGap := f.makeBackground(height)

	gap += extraGap
	x := int((FilmWidth - imageWidth) / 2)
	y := int(gap / 2)
	for _, im := range images {
		targetRect := image.Rect(x, y,
			x+int(imageWidth), y+im.Bounds().Dy())
		draw.Draw(film, targetRect, im, image.Pt(0, 0), draw.Over)
		y += im.Bounds().Dy() + int(gap)
	}

	f.Made = time.Now()
	f.Image = film
	return nil
}

func (stc *Stc) FilmHandler(w http.ResponseWriter, r *http.Request) {
	if time.Since(stc.Film.Made).Seconds() > FilmTimeout {
		err := stc.Film.Update()
		if err != nil {
			log.Printf("failed to update film: %v", err)
		}
	}

	buffer := new(bytes.Buffer)
	err := jpeg.Encode(buffer, stc.Film.Image, nil)
	if err != nil {
		log.Printf("failed to encode film: %v", err)
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	_, err = w.Write(buffer.Bytes())
	if err != nil {
		log.Printf("failed to output film: %v", err)
	}
}
