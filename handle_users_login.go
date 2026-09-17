package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/d4l4-33/chirpy/internal/auth"
	"github.com/d4l4-33/chirpy/internal/database"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding json: %s", err))
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error retreving user: %s", err))
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	tokenDuration, err := time.ParseDuration("1h")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error parsing duration: %s", err))
		return
	}
	token, err := auth.MakeJWT(user.ID, cfg.secret, tokenDuration)
	if err != nil || token == "" {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error creating token: %s", err))
		return
	}

	refreshToken, err := cfg.dbQueries.CreateToken(r.Context(), database.CreateTokenParams{
		Token:     auth.MakeRefreshToken(),
		UserID:    user.ID,
		ExpiresAt: time.Now().AddDate(0, 0, 60),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating refresh token: %s", err))
		return
	}

	err = respondWithJson(w, http.StatusOK, User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken.Token,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

}
