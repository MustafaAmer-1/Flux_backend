package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MustafaAmer-1/Flux/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	decoder := json.NewDecoder(r.Body)

	params := struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't create a feed: %v", err))
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't follow a feed: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, struct {
		feed        Feed
		feed_follow FeedFollow
	}{
		feed:        databaseFeedToFeed(feed),
		feed_follow: databaseFeedFollowToFeedFollow(feedFollow),
	})
}

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get all feeds")
		return
	}
	respondWithJSON(w, http.StatusOK, databaseFeedsToFeeds(feeds))
}
