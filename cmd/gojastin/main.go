package main

import (
	"log"
	"net/http"

	"miikka.xyz/gojastin/server"
)

var buildTime = "This will change when compiled with Makefile"

func main() {
	// How inject string other than main package at compile time?
	s := server.New(buildTime)
	// Delete inactive visitors in background
	go s.CleanVisitors()

	http.HandleFunc("/", s.Router)
	log.Println("Started server on :3030...")
	if err := http.ListenAndServe("0.0.0.0:3030", nil); err != nil {
		log.Fatal(err)
	}
}
