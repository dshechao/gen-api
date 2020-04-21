package main

import (
	"fmt"
	"github.com/dshechao/gen-api/gen"
	"github.com/dshechao/gen-api/middleware"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
	"time"
)

func main() {
	gen.Init(&gen.Config{On: true, DocTitle: "Negroni-gorilla", DocPath: "apidoc.html", BaseUrls: map[string]string{"Production": "", "Staging": ""}})

	router := mux.NewRouter()

	router.HandleFunc("/", middleware.HandleFunc(handler))
	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":5000")
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, time.Now().String())
}
