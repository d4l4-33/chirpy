package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/d4l4-33/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error retrieving header token: %s", err))
		return
	}
	user, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Unauthoarized: %s", err))
		return
	}

	tokenDuration, err := time.ParseDuration("1h")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Unable to parse duration: %s", err))
		return
	}
	accesToken, err := auth.MakeJWT(user.ID, cfg.secret, tokenDuration)

	type parameters struct {
		Token string `json:"token"`
	}

	err = respondWithJson(w, http.StatusOK, parameters{
		Token: accesToken,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error retrieving header token: %s", err))
		return
	}

	err = cfg.dbQueries.RevokeRefreshToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error revoking token: %s", err))
		return
	}

	err = respondWithJson(w, http.StatusNoContent, response{
		Body:  "Refresh token revokded",
		Error: "",
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

}
