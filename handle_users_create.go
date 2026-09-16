package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/d4l4-33/chirpy/internal/auth"
	"github.com/d4l4-33/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token,omitempty"`
}

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request")
		return
	}

	passwordHash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error hashing password: %s", err))
		return
	}

	dbUser, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		HashedPassword: passwordHash,
		Email:          params.Email,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating user: %s", err))
		return
	}

	respondWithJson(w, 201, User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	})
}
