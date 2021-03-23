package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/billglover/albums/internal/albumstore"

	"github.com/gorilla/mux"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	router := mux.NewRouter()
	router.StrictSlash(true)
	server := NewAlbumServer()

	server.store.LoadSamples()

	router.HandleFunc("/album/{id:[0-9]+}/", server.getAlbumHandler).Methods("GET")

	fmt.Fprintln(os.Stdout, "Listenting on localhost:"+os.Getenv("PORT"))
	return http.ListenAndServe("localhost:"+os.Getenv("PORT"), router)
}

type albumServer struct {
	store *albumstore.AlbumStore
}

func NewAlbumServer() *albumServer {
	store := albumstore.New()
	return &albumServer{store: store}
}

// renderJSON renders 'v' as JSON and writes it as a response into w.
func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (as *albumServer) getAlbumHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get album at %s\n", req.URL.Path)

	// Here and elsewhere, not checking error of Atoi because the router only
	// matches the [0-9]+ regex.
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	task, err := as.store.GetAlbum(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, task)
}
