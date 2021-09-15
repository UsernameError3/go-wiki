package main

// Followed Example Project Guide: https://golang.org/doc/articles/wiki/

// Define Imports
import (
	"fmt"
	"io/ioutil"
)

// Define Data Structures
// Slices are similar to arrays but more flexible and more efficient. Reference: https://go.dev/blog/slices-intro
type Page struct {
	Title string
	Body  []byte
}

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

// Main Event Loop
func main() {
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample page.")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))
}
