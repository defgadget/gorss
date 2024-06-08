package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/defgadget/gorss/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	dbURL := os.Getenv("SQL_CONN")

	mux := http.NewServeMux()
	addr := fmt.Sprintf("%s:%s", "localhost", port)
	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}
	dbQueries := database.New(db)
	api := NewAPIConfig(dbQueries)

	mux.HandleFunc("GET /v1/healthz", healthz)
	mux.HandleFunc("GET /v1/err", alwaysErrors)

	mux.HandleFunc("GET /v1/users", api.middlewareAuth(GetUser))
	mux.HandleFunc("POST /v1/users", api.CreateUser)

	mux.HandleFunc("GET /v1/feeds", api.GetAllFeeds)
	mux.HandleFunc("POST /v1/feeds", api.middlewareAuth(api.CreateFeed))

	panic(server.ListenAndServe())
}
