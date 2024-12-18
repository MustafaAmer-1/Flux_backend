package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/MustafaAmer-1/Flux/internal/database"
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
		log.Println("Error marking feed as fetched", err)
		return
	}
	rss_feed, err := fetchFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed", err)
		return
	}

	log.Println("fetched feed with title:", rss_feed.Channel.Title)
}
