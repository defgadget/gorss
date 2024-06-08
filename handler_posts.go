package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/defgadget/gorss/internal/database"
)

func (cfg *apiConfig) GetPostsByUser(rw http.ResponseWriter, req *http.Request, user database.User) {
	limit, err := strconv.Atoi(req.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	posts, err := cfg.DB.GetPostsByUser(req.Context(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		log.Printf("failed to get posts by user: %v", err)
		respondWithError(rw, http.StatusInternalServerError, "ERROR")
		return
	}

	respondWithJSON(rw, http.StatusOK, posts)
}
