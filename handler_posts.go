package main

import (
	"net/http"
	"strconv"

	"github.com/MustafaAmer-1/Flux/internal/database"
)

func (apiCfg *apiConfig) handlerGetPosts(w http.ResponseWriter, r *http.Request, user database.User) {
	limit := 10 // default limit
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
		limit = l
	}
	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get posts")
		return
	}
	respondWithJSON(w, http.StatusOK, databasePostsToPosts(posts))
}
