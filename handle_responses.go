package main

import (
	"encoding/json"
	"net/http"
)

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
