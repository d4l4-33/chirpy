package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	godotenv.Load(".env")
	if cfg.platform != os.Getenv("PLATFORM") {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment"))
		return
	}

	cfg.fileserverHits.Store(0)
	err := cfg.dbQueries.ResetUsers(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to reset the database " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and database reset to intial state."))
}
