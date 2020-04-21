package main

import (
	"fmt"
	"github.com/dshechao/gen-api/gen"
	"github.com/dshechao/gen-api/middleware"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, time.Now().String())
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = ioutil.ReadAll(r.Body)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("test", "tesasasdasd")
	fmt.Fprintf(w, time.Now().String())
}

func main() {
	gen.Init(&gen.Config{On: true, DocTitle: "Gorilla Mux", DocPath: "apidoc.html", BaseUrls: map[string]string{"Production": "", "Staging": ""}})
	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	r.HandleFunc("/testing", postHandler).Methods("POST")
	http.ListenAndServe(":8080", middleware.Handle(r))
}
