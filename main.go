package main

import (
	"gohtmx/handler"
	"log"
	"net/http"
)

const (
	Addr = ":8000"
)

func main() {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    Addr,
		Handler: mux,
	}
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", handler.HomepageRoute)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Keke crash")
	}
}
