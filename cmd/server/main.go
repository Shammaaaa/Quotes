package main

import (
	"Quotes/internal/handlers"
	"Quotes/internal/storage"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	store := storage.New()
	handler := handlers.New(store)

	handler.Routes(r)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
