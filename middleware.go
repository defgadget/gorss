package main

import (
	"context"
	"net/http"

	"github.com/defgadget/gorss/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		apiKey := getAuthorizationHeader(req, "ApiKey")
		ctx := context.Background()
		user, err := cfg.DB.GetUserByAPIKey(ctx, apiKey)
		if err != nil {
			respondWithError(rw, http.StatusInternalServerError, "Error")
			return
		}
		handler(rw, req, user)
	}
}
