package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/MustafaAmer-1/Flux/internal/database"
	"github.com/google/uuid"
)

func startScraping(db *database.Queries, concurrency int32, timeBetweenRequest time.Duration) {
	log.Printf("Scraping on %v goroutines every %v duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), concurrency)
		if err != nil {
			log.Println("error fetching feeds", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go func(feed database.Feed) {
				defer wg.Done()
				scrapeFeed(db, feed)
			}(feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched", feed.Name, err)
		return
	}
	rss_feed, err := fetchFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed", feed.Name, err)
		return
	}

	created_posts := 0
	for _, item := range rss_feed.Channel.Items {
		published_at := sql.NullTime{}
		if time, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			published_at.Time = time
			published_at.Valid = true
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: published_at,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
		created_posts += 1
	}

	log.Printf("Feed %s collected, %v posts found, %v new posts created", feed.Name, len(rss_feed.Channel.Items), created_posts)
}
