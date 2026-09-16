package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/d4l4-33/chirpy/internal/auth"
	"github.com/d4l4-33/chirpy/internal/database"
)

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token available")
		return
	}
	//fmt.Println(tokenString)

	userID, err := auth.ValidateJWT(tokenString, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Unauthorized: %s", err))
		return
	}

	type parameters struct {
		Body        string `json:"body"`
		CleanedBody string `json:"cleaned_body,omitempty"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
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
		UserID: userID,
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
