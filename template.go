package main

import (
	"fmt"
	"html/template"
	"io"
)

const (
	stylesheetTemplate = `
@import url(http://fonts.googleapis.com/css?family=Droid+Sans:400,700);
@import url(http://fonts.googleapis.com/css?family=Special+Elite);
body {
  background: white;
  color: black;
  font-family: "Droid Sans", "dejavu sans", "arial", sans-serif;
  margin-left: 5%;
  margin-right: 5%; }

img {
  border: 0; }

a:link {
  background: transparent;
  color: #0000ff;
  text-decoration: none; }

a:visited {
  background: transparent;
  color: #802020;
  text-decoration: none; }

a:hover {
  text-decoration: underline; }

a.img:link {
  border-bottom: solid white; }

a.img:visited {
  border-bottom: solid white; }

a.img:hover {
  border-bottom: solid black; }

h2 {
  font-family: "Special Elite", "courier", "mono", monospace;
  font-weight: bold;
  font-size: 170%;
  text-align: center; }

h3 {
  font-family: "Special Elite", "courier", "mono", monospace;
  font-weight: bold;
  font-size: 140%; }

h4 {
  font-family: "Special Elite", "courier", "mono", monospace;
  font-weight: bold; }

dt {
  font-family: "Special Elite", "courier", "mono", monospace;
  font-weight: bold; }

dd {
  margin-bottom: 1em; }

span.error {
  font-size: 90%;
  font-weight: bold;
  color: red; }

header.banner {
  margin: 0 auto;
  width: 563px;
  margin-bottom: 50px; }
  header.banner a:hover {
    text-decoration: none; }
  header.banner nav {
    padding-top: 4px; }
    header.banner nav img {
      margin-left: 14px;
      margin-right: 14px; }
    header.banner nav span.social img {
      margin-left: 0px;
      margin-right: 0px; }

footer {
  clear: both; }

section.feature p.image {
  float: left;
  margin-right: 1em;
  margin-bottom: 1em; }
section.feature p.information {
  float: right; }
section.feature p.links {
  clear: both;
  text-align: center;
  font-size: 80%; }
section.feature section.appearances {
  clear: both; }

section.computer p.image {
  float: left;
  margin-right: 1em;
  margin-bottom: 1em; }
section.computer p.information {
  float: right; }
section.computer p.links {
  clear: both;
  text-align: center;
  font-size: 80%; }
section.computer section.appearances {
  clear: both; }

.appearance p.image {
  float: left;
  margin-right: 1em;
  margin-bottom: 1em; }
.appearance p.comment {
  float: right;
  font-weight: bold; }

article.appearance p.stars {
  clear: both; }

section.edgefilm {
  background-image: url("film.jpg");
  background-repeat: repeat-y;
  padding-left: 160px; }

div.film {
  clear: both;
  text-align: center;
  margin: 0 auto;
  width: 528px; }
  div.film div {
    text-align: center;
    background-image: url("/img/film_bg.png");
    background-repeat: repeat-y; }
    div.film div img {
      border: black 2px;
      padding-top: 8px;
      padding-bottom: 8px; }

article.introduction hr {
  width: 60%;
  clear: both; }
article.introduction p.signature {
  float: right; }

section.atoz {
  text-align: center;
  font-size: 80%; }

article.sublist {
  font-size: 80%;
  margin-bottom: 0.5em;
  margin-left: 2em; }

section.help p.image {
  float: left;
  margin-right: 1em;
  margin-bottom: 1em; }
section.help p.information {
  float: right; }
section.help p.links {
  clear: both;
  text-align: center;
  font-size: 80%; }
section.help section.appearances {
  clear: both; }

article.help {
  padding-top: 2em;
  clear: both; }
  article.help img {
    float: left;
    padding-right: 1em; }

article.helped {
  padding-top: 2em;
  clear: both; }
  article.helped img {
    float: left;
    padding-right: 1em; }
  article.helped img {
    width: 200px; }

section.comments {
  clear: both; }
  section.comments article.comment {
    clear: both; }
    section.comments article.comment span.date {
      font-size: 80%; }
  section.comments article.commentform {
    text-align: center;
    float: left; }
    section.comments article.commentform p.approval {
      background: #ccffcc;
      color: #000000;
      border: thin solid #000000;
      text-align: center;
      margin-left: 10%;
      margin-right: 10%;
      margin-top: 1em;
      margin-bottom: 1em;
      padding-top: 1em;
      padding-bottom: 1em; }

section.footerbuttonbox {
  float: left;
  padding-right: 10px; }
  section.footerbuttonbox a:hover {
    text-decoration: none; }

.rssbutton {
  width: 100%;
  text-align: center; }
  .rssbutton a:hover {
    text-decoration: none; }

a.button:hover {
  text-decoration: none; }
	`
	layoutTemplate = `
<html>
<head>
  <title>{{.PageTitle}}</title>
  <!--[if lt IE 9]>
      <script src="http://html5shiv.googlecode.com/svn/trunk/html5.js"></script>
    <![endif]-->
  <meta content='Computers in movies and television shows' name='description'>
  <link href='/stylesheet.css' rel='STYLESHEET' type='text/css'>
  <link href='/favicon.ico' rel='shortcut icon' type='image/x-icon'>
  <script>!function(d,s,id){var js,fjs=d.getElementsByTagName(s)[0];if(!d.getElementById(id)){js=d.createElement(s);js.id=id;js.src="//platform.twitter.com/widgets.js";fjs.parentNode.insertBefore(js,fjs);}}(document,"script","twitter-wjs");</script>
</head>
<body>
<header class='banner'>
    <h1>
      <a href='/'>
        <img alt='Starring the Computer - computers in movies and television' src='/img/banner.png'>
      </a>
    </h1>
    <nav>
      <a class='img' href='/news.html'>
        <img alt='[News]' src='/img/news.png'>
      </a>
      <a class='img' href='/features.html'>
        <img alt='[Movies/TV]' src='/img/movies.png'>
      </a>
      <a class='img' href='/computers.html'>
        <img alt='[Computers]' src='/img/computers.png'>
      </a>
      <a class='img' href='/help.html'>
        <img alt='[Help!]' src='/img/help.png'>
      </a>
      <span class='social'>
        <a class='button' href='http://www.facebook.com/pages/Starring-the-Computer/25393451688'>
          <img alt='[Facebook]' src='/img/facebook.png'>
        </a>
        <a class='button' href='https://twitter.com/StarringTheComp'>
          <img alt='[Twitter]' src='/img/twitter.png'>
        </a>
      </span>
    </nav>
  </header>
</body>
</html>
{{template "content" .}}
  <footer>
    <hr>
    <section class='footerbuttonbox'>
      <a href='http://www.facebook.com/pages/Starring-the-Computer/25393451688'>
        <img alt='[Facebook]' src='/img/facebook.png'>
      </a>
      <a href='https://twitter.com/StarringTheComp'>
        <img alt='[Twitter]' src='/img/twitter.png'>
      </a>
    </section>
    <address>
      <a href='mailto:james@starringthecomputer.com'>
        james@starringthecomputer.com
      </a>
    </address>
  </footer>
</body>
{{define "miniindex"}}
  {{if gt (len .) 1}}
    <p class='links'>
    {{range .}}
      <a href='#{{.Subject.Identity}}'>
        &bull;{{.Subject.Name}}
      </a>
    {{end}}
    </p>
  {{end}}
{{end}}
{{define "star"}}
  <img alt='{{.LabelAlt}}' src='/img/{{.LabelImage}}'>
  <img alt='{{.StarsAlt}}' src='/img/{{.StarsImage}}'>
  {{if .Text}}
    <br>
    {{.Text.Format}}
    <br>
  {{end}}
{{end}}
{{define "stars"}}
        <p class='stars'>
	  {{template "star" .ImportanceInfo}}
	  <br>
	  {{template "star" .RealismInfo}}
	  <br>
	  {{template "star" .VisibilityInfo}}
        </p>
{{end}}
{{define "apppearancelink"}}
        <p class='comment'>
          <a href='/appearance.html?f={{.Feature.Id}}&amp;c={{.Computer.Id}}'>
            Add a comment
	    {{- if gt (len .Comments) 0 -}}
	      &nbsp;({{len .Comments}})
	    {{- end -}}
	    .
          </a>
        </p>
{{end}}
{{define "newsarticles"}}
  <dl>
  {{range .}}
    <dt>{{.Stamp.Format "2006-01-02 15:04"}} {{.Title}}</dt>
    <dd>{{.Text.Format}}</dd>
  {{end}}
  </dl>
{{end}}
    `
	featureTemplate = `
{{define "content"}}
 <section class='feature'>
    <p class='image'>
      <img alt='' src='/movies/{{.Feature.Image}}'>
    </p>
    <h3>{{.Feature.Name}}</h3>
    <p>{{.Feature.Description.Format}}</p>
    <p class='information'>
      <a class='img' href='{{.Feature.ImdbLink}}'>
        <img alt='[More Information]' src='/img/info.png'>
      </a>
    </p>
    {{template "miniindex" .Appearances}}
    <section class='appearances'>
    {{range .Appearances}}
      <article class='appearance' id='{{.Subject.Identity}}'>
	<hr>
	<p class='image'>
          <img alt='' src='/computers/{{.Computer.Image}}'>
        </p>
        <h3>
          <a href='/computer.html?c={{.Computer.Id}}'>{{.Computer.Name}}</a>
        </h3>
        <p>{{.Description.Format}}</p>
	{{template "apppearancelink" .}}
	{{template "stars" .}}
	<div class='film'>
          <div>
	    {{range .Images}}
            <img alt='' src='/snapshots/{{.}}'>
	    {{end}}
          </div>
        </div>
      </article>
    {{end}}
    </section>
  </section>
{{end}}
    `
	computerTemplate = `
{{define "content"}}
 <section class='computer'>
    <p class='image'>
      <img alt='' src='/computers/{{.Computer.Image}}'>
    </p>
    <h3>{{.Computer.Name}}</h3>
    <p>{{.Computer.Description.Format}}</p>
    <p class='information'>
      <a class='img' href='{{.Computer.InfoLink}}'>
        <img alt='[More Information]' src='/img/info.png'>
      </a>
    </p>
    {{template "miniindex" .Appearances}}
    <section class='appearances'>
    {{range .Appearances}}
      <article class='appearance' id='{{.Subject.Identity}}'>
	<hr>
	<p class='image'>
          <img alt='' src='/movies/{{.Feature.Image}}'>
        </p>
        <h3>
          <a href='/feature.html?f={{.Feature.Id}}'>{{.Feature.Name}}</a>
        </h3>
        <p>{{.Description.Format}}</p>
	{{template "apppearancelink" .}}
	{{template "stars" .}}
	<div class='film'>
          <div>
	    {{range .Images}}
	      <img alt='' src='/snapshots/{{.}}'>
	    {{end}}
          </div>
        </div>
      </article>
    {{end}}
    </section>
  </section>
{{end}}
    `
	indexTemplate = `
{{define "content"}}
  <section class='edgefilm'>
    <section class='atoz'>
	{{range $index, $element := .Index.Indices}}
	  {{if gt $index 0}}|{{end}}
	  {{if .Link}} 
	    <a href="#{{.Link}}">{{.Name}}</a>
	  {{else}}
	    {{.Name}}
	  {{end}}
	{{end}}
	{{if .Index.AltName}}
	  <br>
	  <a href='{{.Index.AltLink}}'>sort by {{.Index.AltName}}</a>
	{{end}}
    </section>
    {{range $index, $items := .Index.Entries}}
      <h3 id='{{$index}}'>{{$index}}</h3>
      <ul>
        {{range $items}}
	  <li>
	    <a href='{{.Link}}'>{{.Name}}</a>
	    <article class='sublist'>
	      {{.Things}}
	    </article>
	  </li>
	{{end}}
      </ul>
    {{end}}
  </section>
{{end}}
    `
	introTemplate = `
{{define "content"}}
  <section class='edgefilm'>
  <article class='introduction'>
    <p>
    <b>Starring the Computer</b> is a website dedicated to the use of
    computers in film and television. Each appearance is catalogued and
    rated on its importance (ie. how important it is to the plot), realism
    (how close its appearance and capabilities are to the real thing) and
    visibility (how good a look does one get of it). Fictional computers
    don't count (unless they are built out of bits of real computer), so no
    HAL9000 - sorry.
    </p>
    <p>
    Please <a href="help.html">let me know</a> if you spot any mistakes, or
    have any tips about films not mentioned here that feature computers.
    </p>
    <p class='signature'>
    <a href='http://www.jfc.org.uk'>James Carter</a>
    &lt;<a href='mailto:james@starringthecomputer.com'>james@starringthecomputer.com</a>&gt;
    </p>
    <hr>
  </article>
  <section>
    <h2>News</h2>
    {{template "newsarticles" .News}}
  </section>
{{end}}
    `
	helpTemplate = `
{{define "content"}}
  <section class='edgefilm'>
  <section>
      <h2>Ways to help</h2>
      <p>
        I'd very much appreciate your helping in making
        <a href='/'>starringthecomputer.com</a>
        the best it can be.
      </p>
      <h3>Something missing?</h3>
      <p>
        You can help me by letting me know if there is any movie or TV program
        with a computer in it that I don't already know about. Apart from the
        appearances listed here I have a (rather disorganised)
        <a href='/movies.txt'>list of pending features</a>
        that you should check before
        <a href="mailto:james@starringthecomputer.com">mailing me</a>.
      </p>
      <h3>Most wanted</h3>
      <p>
          Can you identify any of the machines in these pictures? If so,
	  please get in touch!
        </p>
  <a class="twitter-timeline" href="https://twitter.com/StarringTheComp">Tweets by StarringTheComp</a> <script async src="//platform.twitter.com/widgets.js" charset="utf-8"></script>
  </section>
{{end}}
    `
	newsTemplate = `
{{define "content"}}
  <section class='edgefilm'>
  <section>
    <h2>News</h2>
    <section class='atoz'>
      {{if .Newer}}
	<a href='/news.html?{{.Newer}}'>Newer articles</a>
      {{else}}
	Newer articles
      {{end}}
      |
      {{if .Older}}
	<a href='/news.html?{{.Older}}'>Older articles</a>
      {{else}}
	Older articles
      {{end}}
    </section>
    {{template "newsarticles" .News}}
  </section>
  </section>
{{end}}
    `
)

type ComputerTemplateData struct {
	PageTitle   string
	Computer    *Computer
	Appearances []Appearance
}

type FeatureTemplateData struct {
	PageTitle   string
	Feature     *Feature
	Appearances []Appearance
}

type NewsTemplateData struct {
	PageTitle string
	LinkNewer string
	LinkOlder string
	News      []News
}

func (n NewsTemplateData) Newer() template.URL {
	return template.URL(n.LinkNewer)
}

func (n NewsTemplateData) Older() template.URL {
	return template.URL(n.LinkOlder)
}

type IndexTemplate struct {
	PageTitle string
	Index     *Index
}

type Templates map[string]*template.Template

func withLayout(t string) string {
	return layoutTemplate + t
}

func MakeTemplates() (*Templates, error) {
	result := make(Templates)
	for name, tmpl := range map[string]string{
		"index":      withLayout(indexTemplate),
		"feature":    withLayout(featureTemplate),
		"computer":   withLayout(computerTemplate),
		"intro":      withLayout(introTemplate),
		"help":       withLayout(helpTemplate),
		"news":       withLayout(newsTemplate),
		"stylesheet": stylesheetTemplate,
	} {
		t, err := template.New(name).Parse(tmpl)
		if err != nil {
			return &result, err
		}
		result[name] = t
	}
	return &result, nil
}

func (t *Templates) Exec(name string, wr io.Writer, data interface{}) error {
	tmpl, ok := (*t)[name]
	if !ok {
		return fmt.Errorf("no such template %s", name)
	}
	return tmpl.Execute(wr, data)
}
