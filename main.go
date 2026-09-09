package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func handleFileServer(filepathRoot string) http.Handler {

	return http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handleFileServer(filepathRoot)))

	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /api/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("POST /api/reset", apiCfg.handleReset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	err := http.ListenAndServe(srv.Addr, srv.Handler)
	if err != nil {
		fmt.Printf("Error starting server: %s", err)
	}

}
