package main

import (
	"fmt"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.HandleFunc("/healthz", handleReadiness)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	err := http.ListenAndServe(srv.Addr, srv.Handler)
	if err != nil {
		fmt.Printf("Error starting server: %s", err)
	}

}

func handleReadiness(w http.ResponseWriter, _ *http.Request) {

	w.WriteHeader(200)
	w.Write([]byte("OK"))

}
