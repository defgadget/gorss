package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/defgadget/gorss/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) CreateFeed(rw http.ResponseWriter, req *http.Request, user database.User) {
	type requestDTO struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(req.Body)
	reqDTO := requestDTO{}
	err := decoder.Decode(&reqDTO)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "ERROR")
		return
	}

	now := time.Now().UTC()
	feed, err := cfg.DB.CreateFeed(req.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      reqDTO.Name,
		Url:       reqDTO.URL,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "ERROR")
		return
	}
	ff, err := cfg.DB.CreateFeedFollow(req.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: now,
		UpdatedAt: now,
	})

	responseDTO := struct {
		Feed       database.Feed      `json:"feed"`
		FeedFollow database.UsersFeed `json:"feed_follow"`
	}{Feed: feed, FeedFollow: ff}

	respondWithJSON(rw, http.StatusCreated, responseDTO)
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
