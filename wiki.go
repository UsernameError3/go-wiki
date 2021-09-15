package main

// Followed Example Project Guide: https://golang.org/doc/articles/wiki/

// Define Imports
import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"text/template"
)

// Define Variables
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// Define Data Structures
// Slices are similar to arrays but more flexible and more efficient. Reference: https://go.dev/blog/slices-intro
type Page struct {
	Title string
	Body  []byte
}

// Define Functions
// Save Method for Persistent Storage
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// Loading Pages and Catching Errors
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// References the HTML via Templates rather than hardcoding for better readability.
// Uses the ExecuteTemplate() method to effectively cache templates at init, then calling specific templates as necessary.
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Validates the Path to prevent users supplying arbitrary ptahs to be read/written on the server.
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Title")
	}
	return m[2], nil
}

// Use Function Literals and Closures to wrap each handler to remove validation redundancy
/*
Notes from guide:
The closure returned by makeHandler is a function that takes an http.ResponseWriter and http.Request (in other words, an http.HandlerFunc).
The closure extracts the title from the request path, and validates it with the validPath regexp.
If the title is invalid, an error will be written to the ResponseWriter using the http.NotFound function.
If the title is valid, the enclosed handler function fn will be called with the ResponseWriter, Request, and title as arguments.
*/
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

// Allows users to view a Wiki Page by handling URL's with '/view/'
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// Allows users to edit a Wiki Page by handling URL's with '/edit/'
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

// Allows users to save a Wiki Page after using the edit function, and redirecting back to the View Handler with '/view/'
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

// Main Event Loop
func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
