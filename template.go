package main

import (
	"fmt"
	"html/template"
	"io"
	"strings"
	"time"
)

const (
	// Time formats ---------------------------------------------------
	RssTimeFormat   = "Mon, 02 Jan 2006 15:04:05 -0700"
	IndexTimeFormat = "2006-01-02 15:04"

	// RSS ------------------------------------------------------------
	rssTemplate = `
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <atom:link href="http://www.starringthecomputer.com/stc.rss" rel="self" type="application/rss+xml" />
    <title>Starring the Computer</title>
    <link>http://www.starringthecomputer.com</link>
    <description>Starring the Computer is a website dedicated to the use of computer in movies and television.</description>
    <language>en</language>
    <copyright>Copyright 2007-{{.Now.Year}} James Carter</copyright>
    <lastBuildDate>{{rssTime .Now}}</lastBuildDate>
    <image>
      <url>http://www.starringthecomputer.com/img/starringthecomputer.png</url>
      <title>Starring the Computer</title>
      <link>http://www.starringthecomputer.com</link>
    </image>
    {{range .News}}
    <item>
      <title>{{.Title}}</title>
      <link>http://www.starringthecomputer.com/newsitem.html?i={{rssIndexURL .Stamp}}</link>
      <description>{{html .Text.FormatFullUrl}}</description>
      <pubDate>{{rssTime .Stamp}}</pubDate>
      <guid isPermaLink="true">http://www.starringthecomputer.com/newsitem.html?i={{rssIndexURL .Stamp}}</guid>
    </item>
    {{end}}
  </channel>
</rss>
    `

	// CSS ------------------------------------------------------------
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
    header.banner nav.menu {
      margin-top: 20px;
      margin-left: 48px; }
    header.banner nav.menu img {
      margin-left: 14px;
      margin-right: 14px; }
    header.banner nav.social {
      margin-top: -64px;
      margin-left: 352px;
      width: 150px; }
    header.banner nav.social img {
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

	// LAYOUT ---------------------------------------------------------
	layoutTemplate = `
<html>
<head>
  <title>{{.PageTitle}}</title>
  {{ifIElt9}}
      <script src="http://html5shiv.googlecode.com/svn/trunk/html5.js"></script>
  {{endif}}
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
    <nav class='social'>
      <a class='button' href='http://www.facebook.com/pages/Starring-the-Computer/25393451688'>
        <img alt='[Facebook]' src='/img/facebook.png'>
      </a>
      <a class='button' href='https://twitter.com/StarringTheComp'>
        <img alt='[Twitter]' src='/img/twitter.png'>
      </a>
      <a class='button' href='/stc.rss'>
        <img alt='[RSS]' src='/img/rss.png'>
      </a>
      <a class='button' href='https://ko-fi.com/starringthecomputer'>
        <img alt='[RSS]' src='/img/kofi.png'>
      </a>
    </nav>
    <nav class='menu'>
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
    </nav>
  </header>
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
      <a href='/stc.rss'>
        <img alt='[RSS]' src='/img/rss.png'>
      </a>
    </section>
    <address>
      <a href='mailto:james@starringthecomputer.com'>
        james@starringthecomputer.com
      </a>
    </address>
  </footer>
</body>
</html>
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
{{define "appearancelink"}}
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
    <dt>{{indexTime .Stamp}} {{.Title}}</dt>
    <dd>{{.Text.Format}}</dd>
  {{end}}
  </dl>
{{end}}
{{define "film"}}
  <div class='film'>
    <div>
      {{range .Images}}
      <img alt='' src='/snapshots/{{.}}'>
      {{end}}
    </div>
  </div>
{{end}}
    `

	// HTML -----------------------------------------------------------
	appearanceTemplate = `
{{define "content"}}
  <section class='appearance'>
    <p class='image'>
      <img alt='' src='/computers/{{.Appearance.Computer.Image}}'>
      <br>
      <img alt='' src='/movies/{{.Appearance.Feature.Image}}'>
    </p>
    <h2>{{.Appearance.Computer.Name}} in {{.Appearance.Feature.Name}}</h2>
    <p>{{.Appearance.Description.Format}}</p>
    {{template "stars" .Appearance}}
    <section class='comments'>
      <h3><img alt='Comments:' href='/img/comments.png'></h3>
      <article class='commentform'>
      {{if .Form.Created}}
        <p class='approval'>
	  Your comment has been submitted for approval and will appear here
	  shortly. Thanks!
	</p>
      {{else}}
        <form action='/appearance.html' method='post'>
	  <p>
	    <input name='c' type='hidden' value='{{.Appearance.Computer.Id}}'>
	    <input name='f' type='hidden' value='{{.Appearance.Feature.Id}}'>
	    {{.Form.Label "Name"}}
	    <input name='n' type='text' value='{{.Form.Name}}'><br><br>
	    {{.Form.Label "Comment"}}<br>
	    <textarea cols='60' name='t' rows='6'>{{.Form.Comment}}</textarea><br><br>
	    {{.Form.Label "Year of feature (shown above)"}}
	    <input name='y' type='text' value='{{.Form.Year}}'><br><br>
	    <input name='post' type='submit' value='Post Comment'>
	  </p>
	</form>
      {{end}}
      </article>
      {{range .Appearance.Comments}}
        <article class='comment'>
	  <hr>
	  <h4>{{.Name}}</h4>
	  <p>
	    {{.Text}}
	    <br><span class='date'>{{indexTime .Stamp}}</span>
	  </p>
        </article>
      {{end}}
    </section>
    {{template "film" .Appearance}}
  </section>
{{end}}
    `
	featureTemplate = `
{{define "content"}}
 <section class='feature'>
    <p class='image'>
      <img alt='' src='/movies/{{.Feature.Image}}'>
    </p>
    <h2>{{.Feature.Name}}</h2>
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
	{{template "appearancelink" .}}
	{{template "stars" .}}
	{{template "film" .}}
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
      {{if .Computer.ImageLink.Valid -}}
        <a class='img' href='{{.Computer.ImageLink.String}}'>
      {{- end}}
	<img alt='' src='/computers/{{.Computer.Image}}'>
      {{if .Computer.ImageLink.Valid -}}
        </a>
      {{- end}}
    </p>
    <h2>{{.Computer.Name}}</h2>
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
	{{template "appearancelink" .}}
	{{template "stars" .}}
	{{template "film" .}}
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
    {{range $index, $element := .Index.OrderedEntries}}
      <h3 id='{{.Index}}'>{{.Index}}</h3>
      <ul>
        {{range .Entries}}
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
      <h3>Honourable mentions</h3>
      <p>
        The following films do not appear on the site because I believe
	the computers they feature are mock ups and therefore do not
	qualify. If you have concrete evidence that those machines are
	real, or are made from identifiable parts of other machines then
	do feel free to get in touch!
      </p>
      <dl>
      	<dt>The Desk Set</dt>
	<dd>Although it uses IBM tape drives, this machine is just a mock up.</dd>
      	<dt>Sneakers</dt>
	<dd>The Cray-like machine appears to be a mock up.</dd>
      </dl>
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
{{define "nav"}}
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
{{end}}
{{define "content"}}
  <section class='edgefilm'>
  <section>
    <h2>News</h2>
    {{template "nav" .}}
    {{template "newsarticles" .News}}
    {{template "nav" .}}
  </section>
  </section>
{{end}}
    `
	newsitemTemplate = `
{{define "content"}}
  <section class='edgefilm'>
  <section>
    {{template "newsarticles" .News}}
  </section>
  </section>
{{end}}
    `
	approvecommentTemplate = `
{{define "content"}}
    <h1>Approved comment</h1>
    <p>View it here: <a href='{{.Link}}'>{{.Link}}</a></p>
{{end}}
    `
	denycommentTemplate = `
{{define "content"}}
    <h1>Deny comment</h1>
    <p>View it here: <a href='{{.Link}}'>{{.Link}}</a></p>
    <p>Click here to delete comment: <a href='{{.DelLink}}'>{{.DelLink}}</a></p>
{{end}}
    `
	deletecommentTemplate = `
{{define "content"}}
    <h1>Comment deleted</h1>
{{end}}
    `
)

type AppearanceTemplateData struct {
	PageTitle  string
	Appearance *Appearance
	Form       *CommentForm
}

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

type RssTemplateData struct {
	Now  time.Time
	News []News
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
		"index":          withLayout(indexTemplate),
		"appearance":     withLayout(appearanceTemplate),
		"feature":        withLayout(featureTemplate),
		"computer":       withLayout(computerTemplate),
		"intro":          withLayout(introTemplate),
		"help":           withLayout(helpTemplate),
		"news":           withLayout(newsTemplate),
		"newsitem":       withLayout(newsitemTemplate),
		"approvecomment": withLayout(approvecommentTemplate),
		"denycomment":    withLayout(denycommentTemplate),
		"deletecomment":  withLayout(deletecommentTemplate),
		"stylesheet":     stylesheetTemplate,
		"rss":            rssTemplate,
	} {
		t, err := template.New(name).Funcs(template.FuncMap{
			"indexTime": func(t time.Time) string {
				return t.Format(IndexTimeFormat)
			},
			"rssTime": func(t time.Time) string {
				return t.Format(RssTimeFormat)
			},
			"rssIndexURL": func(t time.Time) template.URL {
				s := t.Format(IndexTimeFormat)
				s = strings.Replace(s, " ", "%20", -1)
				s = strings.Replace(s, "+", "%2B", -1)
				return template.URL(s)
			},
			"ifIElt9": func() template.HTML {
				return template.HTML("<!--[if lt IE 9]>")
			},
			"endif": func() template.HTML {
				return template.HTML("<![endif]-->")
			},
		}).Parse(tmpl)
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
