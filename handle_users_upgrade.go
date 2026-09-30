package main

import (
	"encoding/json"
	"net/http"

	"github.com/primusprag/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleWebhooks(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			User_id uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	requestKey, err := auth.GetAPIKey(r.Header)
	if err != nil || requestKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJson(w, http.StatusNoContent, nil)
		return
	}

	err = cfg.dbQueries.UpgradeUserByID(r.Context(), params.Data.User_id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	err = respondWithJson(w, http.StatusNoContent, nil)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
