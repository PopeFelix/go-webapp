package main

// Based on https://www.sohamkamani.com/blog/2017/09/13/how-to-build-a-web-application-in-golang/
import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/hello", handler).Methods("GET")

	staticDir := http.Dir("./assets/")

	staticHandler := http.StripPrefix("/assets/", http.FileServer(staticDir))

	r.PathPrefix("/assets").Handler(staticHandler).Methods("GET")
	return r
}

func main() {
	r := newRouter()
	http.ListenAndServe(":8080", r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "fnord\n")
}
