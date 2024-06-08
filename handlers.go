package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/defgadget/gorss/internal/database"
	"github.com/google/uuid"
)

type apiConfig struct {
	DB *database.Queries
}

func NewAPIConfig(db *database.Queries) *apiConfig {
	return &apiConfig{DB: db}
}

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

func (cfg *apiConfig) CreateFeed(rw http.ResponseWriter, req *http.Request, user database.User) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "ERROR")
		return
	}

	requestDTO := struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}{}
	err = json.Unmarshal(data, &requestDTO)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "ERROR")
		return
	}

	ctx := context.Background()
	feed, err := cfg.DB.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      requestDTO.Name,
		Url:       requestDTO.URL,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "ERROR")
		return
	}

	respondWithJSON(rw, http.StatusCreated, feed)
}

func (cfg *apiConfig) GetAllFeeds(rw http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	feeds, err := cfg.DB.AllFeeds(ctx)
	if err != nil {
		log.Printf("failed to retrieve feeds: %v", err)
		respondWithError(rw, http.StatusInternalServerError, "error")
		return
	}

	respondWithJSON(rw, http.StatusOK, feeds)
}

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

type userDTO struct {
	Name string `json:"name"`
}
