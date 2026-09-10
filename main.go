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
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleReset)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.handleValidate)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	err := http.ListenAndServe(srv.Addr, srv.Handler)
	if err != nil {
		fmt.Printf("Error starting server: %s", err)
	}

}
