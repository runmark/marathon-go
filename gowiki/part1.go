package main

import (
	"fmt"
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

func main() {
	p1 := &Page{"TestPage", []byte("This is a sample page.")}
	p1.save()

	p2, _ := load("TestPage")
	fmt.Println(p2.Body)
}
