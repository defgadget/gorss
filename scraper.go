package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/araddon/dateparse"
	"github.com/defgadget/gorss/internal/database"
	"github.com/google/uuid"
)

type Scraper struct {
	ConcurrentJobs int
	wg             sync.WaitGroup
	db             *database.Queries
}

func (s *Scraper) Run() {
	ticker := time.NewTicker(60 * time.Second)
	ctx := context.Background()
	for {
		log.Println("Starting to scrape")
		feedsToFetch, err := s.db.GetNextFeedsToFetch(ctx, int32(s.ConcurrentJobs))
		if err != nil {
			log.Printf("failed to fetch next feeds: %v", err)
			break
		}

		for _, feed := range feedsToFetch {
			s.wg.Add(1)
			go func(feed database.Feed, wg *sync.WaitGroup) {
				log.Printf("Processing %s", feed.Url)
				items, err := fetchFeedPosts(feed.Url)
				if err != nil {
					log.Printf("failed to fetch feed posts: %v", err)
					wg.Done()
					return
				}
				s.db.MarkFeedFetched(context.Background(), feed.ID)
				for _, item := range items {
					now := time.Now().UTC()
					pubDate, _ := dateparse.ParseIn(item.PubDate, time.UTC)
					_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
						ID:            uuid.New(),
						CreatedAt:     now,
						UpdatedAt:     now,
						Title:         item.Title,
						Description:   descriptionToSQLNullString(item.Description),
						Url:           item.Link,
						PublishedDate: pubDate,
						FeedID:        feed.ID,
					})
					if err != nil {
						if strings.Contains(err.Error(), "duplicate key") {
							continue
						}
						log.Printf("failed to create post: %v", err)
					}
				}
				wg.Done()
			}(feed, &s.wg)
		}
		log.Println("Waiting...")
		s.wg.Wait()
		log.Println("Finished")
		<-ticker.C
	}
}

func fetchFeedPosts(url string) ([]Item, error) {
	log.Println("Fetching:", url)
	client := http.DefaultClient
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("failed to fetch from (%s): %v", url, err)
		return nil, err
	}
	decoder := xml.NewDecoder(resp.Body)
	rss := Rss{}
	err = decoder.Decode(&rss)
	if err != nil {
		log.Printf("failed to unmarshal XML: %v", err)
		return nil, err
	}
	return rss.Channel.Items, nil
}

func descriptionToSQLNullString(desc string) sql.NullString {
	n := sql.NullString{}
	if desc == "" {
		n.Valid = false
		return n
	}
	n.Valid = true
	n.String = desc

	return n
}

func printItemsList(items []Item) {
	for _, item := range items {
		log.Println("Title:", item.Title)
	}
}
