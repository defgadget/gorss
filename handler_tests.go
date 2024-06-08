package main

import "net/http"

func healthz(w http.ResponseWriter, r *http.Request) {
	_ = r
	status := struct {
		Status string `json:"status"`
	}{Status: "ok"}
	respondWithJSON(w, http.StatusOK, status)
}

func alwaysErrors(w http.ResponseWriter, r *http.Request) {
	_ = r
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
