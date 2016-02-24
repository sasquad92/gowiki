//
// Code from tutorial avaible at <https://golang.org/doc/articles/wiki/>
//  

//
// TODO: (from bottom of tutorial web page)
// 1) Add a handler to make the web root redirect to /view/FrontPage.
// 2) Spruce up the page templates by making them valid HTML and adding some CSS rules.
// 3) Implement inter-page linking by converting instances of [PageName] to
// <a href="/view/PageName">PageName</a>. (hint: you could use regexp.ReplaceAllFunc to do this) 
//

package main

import (
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// The page struct describes how page data will be stored in memory
type Page struct {
	Title string
	Body  []byte
}

// Func unused; see makeHandler
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) { 
	m := validPath.FindStringSubmatch(r.URL.Path)

	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}

	return m[2], nil // The title is second subexpression
}

// Saving pages to data folder
func (p *Page) save() error {
	filename := "data/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// Loading pages from data folder
func loadPage(title string) (*Page, error) {
	filename := "data/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}


func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	filename := tmpl + ".html"

	err := templates.ExecuteTemplate(w, filename, p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)

	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)

	if err != nil {
		p = &Page{Title: title}
	}

	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Here we will etract the page title form the Request,
		// and all the provided handler 'fn'
		m := validPath.FindStringSubmatch(r.URL.Path)

		if m == nil {
			http.NotFound(w, r) // error msg
			return
		}
		fn(w, r, m[2])
	}
}

func main() {

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":8080", nil)

}
