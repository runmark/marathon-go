package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func load(title string) (*Page, error) {
	filename := title + ".txt"

	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &Page{title, body}, nil
}

func handler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "Hi there, I love %s!", r.URL.Path[1:])
}

func renderTemplate(rw http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(rw, p)
}

func viewHandler(rw http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := load(title)

	renderTemplate(rw, "view", p)
	// fmt.Fprintf(rw, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func editHandler(rw http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := load(title)
	if err != nil {
		p = &Page{Title: title}
	}

	renderTemplate(rw, "edit", p)

	// t, _ := template.ParseFiles("edit.html")
	// t.Execute(rw, p)

	// fmt.Fprintf(rw, `
	// 	<h1>Editing %s</h1>
	// 	<form action="/save/%s" method="POST">
	// 		<textarea name="body">%s</textarea><br>
	// 		<input type="submit" value="Save">
	// 	</form>
	// `, p.Title, p.Title, p.Body)
}

func saveHandler(rw http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
