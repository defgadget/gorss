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
	DB  *database.Queries
	ctx context.Context
}

func NewAPIConfig(db *database.Queries, ctx context.Context) *apiConfig {
	return &apiConfig{DB: db, ctx: ctx}
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

	now := time.Now().UTC()
	feed, err := cfg.DB.CreateFeed(cfg.ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      requestDTO.Name,
		Url:       requestDTO.URL,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "ERROR")
		return
	}
	ff, err := cfg.DB.CreateFeedFollow(cfg.ctx, database.CreateFeedFollowParams{
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
	ff, err := cfg.DB.CreateFeedFollow(cfg.ctx, database.CreateFeedFollowParams{
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
	err = cfg.DB.DeleteFeedFollow(cfg.ctx, database.DeleteFeedFollowParams{
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
	ff, err := cfg.DB.AllUserFeedFollows(cfg.ctx, user.ID)
	if err != nil {
		log.Printf("failed to retrieve feed follows: %v", err)
		respondWithError(rw, http.StatusInternalServerError, "ERROR")
		return
	}
	respondWithJSON(rw, http.StatusOK, ff)
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
