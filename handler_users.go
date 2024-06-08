package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/defgadget/gorss/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	userRequest := userDTO{}
	err = json.Unmarshal(data, &userRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	ctx := context.Background()
	user, err := cfg.DB.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      userRequest.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func GetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, user)
}
