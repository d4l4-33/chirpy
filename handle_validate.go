package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type parameters struct {
	Body         string `json:"body,omitempty"`
	Error        string `json:"error,omitempty"`
	Cleaned_body string `json:"cleaned_body,omitempty"`
}

func (cfg *apiConfig) handleValidate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Error decoding request")
		return
	}

	if len(params.Body) >= 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	params.cleanBody()
	respondWithJson(w, 200, params)
}

func (par *parameters) cleanBody() {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	splitBody := strings.Split(par.Body, " ")
	for i, word := range splitBody {
		for _, bWord := range badWords {
			if strings.ToLower(word) == bWord {
				splitBody[i] = "****"
			}
		}
	}
	par.Cleaned_body = strings.Join(splitBody, " ")
}
