/*
	The only reason for this model file is to own the shape of the data modles which being returned from my server,
	not to skick with the ones defined by sqlc,
	Although this is purly optional and we can use the auto generated data modles
*/

package main

import (
	"database/sql"
	"time"

	"github.com/MustafaAmer-1/Flux/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func databaseUserToUser(dbuser database.User) User {
	return User{
		ID:        dbuser.ID,
		CreatedAt: dbuser.CreatedAt,
		UpdatedAt: dbuser.UpdatedAt,
		Name:      dbuser.Name,
		ApiKey:    dbuser.ApiKey,
	}
}

type Feed struct {
	ID            uuid.UUID  `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Name          string     `json:"name"`
	Url           string     `json:"url"`
	UserID        uuid.UUID  `json:"user_id"`
	LastFetchedAt *time.Time `json:"last_fetched_at"`
}

func databaseFeedToFeed(feed database.Feed) Feed {
	return Feed{
		ID:            feed.ID,
		CreatedAt:     feed.CreatedAt,
		UpdatedAt:     feed.UpdatedAt,
		Name:          feed.Name,
		Url:           feed.Url,
		UserID:        feed.UserID,
		LastFetchedAt: nullTimeToTimePtr(feed.LastFetchedAt),
	}
}

func databaseFeedsToFeeds(dbfeeds []database.Feed) []Feed {
	feeds := make([]Feed, len(dbfeeds))
	for i, feed := range dbfeeds {
		feeds[i] = databaseFeedToFeed(feed)
	}
	return feeds
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func databaseFeedFollowToFeedFollow(dbFeedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        dbFeedFollow.ID,
		CreatedAt: dbFeedFollow.CreatedAt,
		UpdatedAt: dbFeedFollow.UpdatedAt,
		UserID:    dbFeedFollow.UserID,
		FeedID:    dbFeedFollow.FeedID,
	}
}

func databaseFeedFollowsToFeedFollows(dbFeedFollows []database.FeedFollow) []FeedFollow {
	feedFollows := make([]FeedFollow, len(dbFeedFollows))
	for i, dbFeedFollow := range dbFeedFollows {
		feedFollows[i] = databaseFeedFollowToFeedFollow(dbFeedFollow)
	}
	return feedFollows
}

type Post struct {
	ID          uuid.UUID  `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Title       string     `json:"title"`
	Url         string     `json:"url"`
	Description *string    `json:"description"`
	PublishedAt *time.Time `json:"published_at"`
	FeedID      uuid.UUID  `json:"feed_id"`
}

func databasePostToPost(dbpost database.Post) Post {
	return Post{
		ID:          dbpost.ID,
		CreatedAt:   dbpost.CreatedAt,
		UpdatedAt:   dbpost.UpdatedAt,
		Title:       dbpost.Title,
		Url:         dbpost.Url,
		Description: nullStringToStringPtr(dbpost.Description),
		PublishedAt: nullTimeToTimePtr(dbpost.PublishedAt),
		FeedID:      dbpost.FeedID,
	}
}

func databasePostsToPosts(dbposts []database.Post) []Post {
	posts := make([]Post, len(dbposts))
	for i, dbpost := range dbposts {
		posts[i] = databasePostToPost(dbpost)
	}
	return posts
}

func nullTimeToTimePtr(time sql.NullTime) *time.Time {
	if time.Valid {
		return &time.Time
	}
	return nil
}

func nullStringToStringPtr(str sql.NullString) *string {
	if str.Valid {
		return &str.String
	}
	return nil
}
