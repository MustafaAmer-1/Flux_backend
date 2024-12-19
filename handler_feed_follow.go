package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MustafaAmer-1/Flux/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerFollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	decoder := json.NewDecoder(r.Body)

	params := struct {
		FeedId uuid.UUID `json:"feed_id"`
	}{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    params.FeedId,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't follow a feedFollow: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedFollowToFeedFollow(feedFollow))
}

func (apiCfg *apiConfig) handlerGetFollowedFeeds(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowed, err := apiCfg.DB.GetFeedsFollowByUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get followed feeds")
		return
	}
	respondWithJSON(w, http.StatusOK, databaseFeedsToFeeds(feedFollowed))
}

func (apiCfg *apiConfig) handlerGetNotFollowedFeeds(w http.ResponseWriter, r *http.Request, user database.User) {
	feeds, err := apiCfg.DB.GetFeedsNotFollowByUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get feeds")
		return
	}
	respondWithJSON(w, http.StatusOK, databaseFeedsToFeeds(feeds))
}

func (apiCfg *apiConfig) handlerUnFollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	feed_id, err := uuid.Parse(chi.URLParam(r, "feedId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid feed id")
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		FeedID: feed_id,
		UserID: user.ID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't unfollow this feed")
		return
	}
	respondWithJSON(w, http.StatusOK, "Feed Unfollowed")
}
