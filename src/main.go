package main

import (
	"log"
	"net/http"
	"weight-tracker/src/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// serve single-page frontend (serve the file directly to avoid redirect loops)
	serveIndex := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "webapp/index.html")
	}
	r.HandleFunc("/", serveIndex).Methods("GET")
	r.HandleFunc("/index.html", serveIndex).Methods("GET")

	// serve static assets
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("webapp/static"))))

	// weight tracking endpoints
	r.HandleFunc("/weights", handlers.ListWeights).Methods("GET")
	r.HandleFunc("/weights", handlers.CreateWeight).Methods("POST")
	r.HandleFunc("/weights/{id}", handlers.GetWeight).Methods("GET")
	r.HandleFunc("/weights/{id}", handlers.UpdateWeight).Methods("PUT")
	r.HandleFunc("/weights/{id}", handlers.DeleteWeight).Methods("DELETE")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
