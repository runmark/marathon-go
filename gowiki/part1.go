package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := dataDir + p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func load(title string) (*Page, error) {
	filename := dataDir + title + ".txt"

	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &Page{title, body}, nil
}

// func handler(rw http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(rw, "Hi there, I love %s!", r.URL.Path[1:])
// }

func renderTemplate(rw http.ResponseWriter, tmpl string, p *Page) {

	err := templates.ExecuteTemplate(rw, tmpl+".html", p)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {

	return func(rw http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(rw, r)
			return
		}
		fn(rw, r, m[2])
	}
}

func rootPathHandler(rw http.ResponseWriter, r *http.Request) {
	rootPath := r.URL.Path
	if rootPath != "/" {
		http.NotFound(rw, r)
		return
	}
	http.Redirect(rw, r, "/view/FrontPage", http.StatusFound)
}

func viewHandler(rw http.ResponseWriter, r *http.Request, title string) {

	p, err := load(title)
	if err != nil {
		http.Redirect(rw, r, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(rw, "view", p)
}

func editHandler(rw http.ResponseWriter, r *http.Request, title string) {

	p, err := load(title)
	if err != nil {
		p = &Page{Title: title}
	}

	renderTemplate(rw, "edit", p)
}

func saveHandler(rw http.ResponseWriter, r *http.Request, title string) {

	body := r.FormValue("body")
	p := &Page{title, []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/view/"+title, http.StatusFound)
}

const tmplDir = "./tmpl/"
const dataDir = "./data/"

var templates = template.Must(template.ParseFiles(tmplDir+"view.html", tmplDir+"edit.html"))
var validPath = regexp.MustCompile("^/(view|edit|save)/([0-9a-zA-Z]+)$")

func main() {
	http.HandleFunc("/", rootPathHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
