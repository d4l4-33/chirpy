package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	dbChirps, err := cfg.dbQueries.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error retrieving chirps: %s", err))
		return
	}

	parsedChirps := []Chirp{}
	if len(dbChirps) != 0 {
		for _, dbchirp := range dbChirps {
			parsedChirps = append(parsedChirps, Chirp{
				ID:        dbchirp.ID,
				CreatedAt: dbchirp.CreatedAt,
				UpdatedAt: dbchirp.UpdatedAt,
				Body:      dbchirp.Body,
				UserID:    dbchirp.UserID,
			})
		}
	}

	err = respondWithJson(w, http.StatusOK, parsedChirps)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

}

func (cfg *apiConfig) handleGetChirpByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	parsedID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing chirpID: %s", err))
		return
	}

	dbChirp, err := cfg.dbQueries.GetChirpByID(r.Context(), parsedID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	err = respondWithJson(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

}
