package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	error := struct {
		Error string `json:"error"`
	}{Error: msg}
	respondWithJSON(w, code, error)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	json, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	// w.Write(json)
	fmt.Fprint(w, string(json))
}

func getAuthorizationHeader(r *http.Request, prefix string) string {
	header := r.Header.Get("Authorization")
	value := strings.TrimPrefix(header, prefix+" ")
	return value
}
