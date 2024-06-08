package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/defgadget/gorss/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) CreateFeedFollow(rw http.ResponseWriter, req *http.Request, user database.User) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "ERROR")
		return
	}
	requestDTO := struct {
		FeedID uuid.UUID `json:"feed_id"`
	}{}
	err = json.Unmarshal(data, &requestDTO)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "ERROR")
		return
	}
	now := time.Now().UTC()
	ff, err := cfg.DB.CreateFeedFollow(req.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    requestDTO.FeedID,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		log.Printf("failed to create feed follow: %v", err)
		respondWithError(rw, http.StatusInternalServerError, "ERROR")
		return
	}

	respondWithJSON(rw, http.StatusCreated, ff)
}

func (cfg *apiConfig) DeleteFollowFeed(rw http.ResponseWriter, req *http.Request, user database.User) {
	path := req.PathValue("feedFollowID")
	ffID, err := uuid.Parse(path)
	if err != nil {
		log.Printf("failed to parse UUID: %v", err)
		respondWithError(rw, http.StatusInternalServerError, "ERROR")
		return
	}
	err = cfg.DB.DeleteFeedFollow(req.Context(), database.DeleteFeedFollowParams{
		ID:     ffID,
		UserID: user.ID,
	})
	if err != nil {
		log.Printf("failed to delete feed follow: %v", err)
		respondWithError(rw, http.StatusNotFound, "failed")
		return
	}

	respondWithJSON(rw, http.StatusNoContent, struct{}{})
}

func (cfg *apiConfig) GetAllFeedFollows(rw http.ResponseWriter, req *http.Request, user database.User) {
	ff, err := cfg.DB.AllUserFeedFollows(req.Context(), user.ID)
	if err != nil {
		log.Printf("failed to retrieve feed follows: %v", err)
		respondWithError(rw, http.StatusInternalServerError, "ERROR")
		return
	}
	respondWithJSON(rw, http.StatusOK, ff)
}

type userDTO struct {
	Name string `json:"name"`
}
