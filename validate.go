package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type parameters struct {
	Body         string `json:"body,omitempty"`
	Error        string `json:"error,omitempty"`
	Cleaned_body string `json:"cleaned_body,omitempty"`
}

func (cfg *apiConfig) handleValidate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
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
	fmt.Println(badWords)
	for i, word := range splitBody {
		fmt.Println(word)
		for _, bWord := range badWords {
			if strings.ToLower(word) == bWord {
				fmt.Printf("bad word found! %s = %s\n", word, bWord)
				splitBody[i] = "****"
			}
		}
	}
	par.Cleaned_body = strings.Join(splitBody, " ")
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJson(w, code, parameters{
		Body:  "",
		Error: msg,
	})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return respondWithError(w, 500, "Error marshaling json")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)

	return nil
}
