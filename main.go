package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/defgadget/gorss/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
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

	mux.HandleFunc("GET /v1/feed_follows", api.middlewareAuth(api.GetAllFeedFollows))
	mux.HandleFunc("POST /v1/feed_follows", api.middlewareAuth(api.CreateFeedFollow))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", api.middlewareAuth(api.DeleteFollowFeed))

	mux.HandleFunc("GET /v1/posts", api.middlewareAuth(api.GetPostsByUser))

	concurrentJobs := 10
	scraper := Scraper{ConcurrentJobs: concurrentJobs, wg: sync.WaitGroup{}, db: dbQueries}
	go scraper.Run()
	panic(server.ListenAndServe())
}

type apiConfig struct {
	DB *database.Queries
}

func NewAPIConfig(db *database.Queries) *apiConfig {
	return &apiConfig{DB: db}
}
