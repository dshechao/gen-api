package main

import (
	"github.com/dshechao/gen-api/gen"
	"github.com/dshechao/gen-api/martiniyaag"
	"github.com/go-martini/martini"
)

func main() {
	gen.Init(&gen.Config{On: true, DocTitle: "Martini", DocPath: "apidoc.html", BaseUrls: map[string]string{"Production": "", "Staging": ""}})
	m := martini.Classic()
	m.Use(martiniyaag.Document)
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Run()
}
