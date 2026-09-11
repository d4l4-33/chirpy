package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/d4l4-33/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type parameters struct {
		Body        string    `json:"body"`
		UserID      uuid.UUID `json:"user_id"`
		CleanedBody string    `json:"cleaned_body,omitempty"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding request: %s", err))
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error creating chirp: %s", err))
		return
	}

	respondWithJson(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func validateChirp(body string) (string, error) {
	if len(body) >= 140 {
		return "", errors.New("Chirp is too long")
	}

	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	splitBody := strings.Split(body, " ")
	for i, word := range splitBody {
		for _, bWord := range badWords {
			if strings.ToLower(word) == bWord {
				splitBody[i] = "****"
			}
		}
	}
	return strings.Join(splitBody, " "), nil

}
