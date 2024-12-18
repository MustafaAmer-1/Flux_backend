package main

import (
	"fmt"
	"net/http"

	"github.com/MustafaAmer-1/Flux/internal/auth"
	"github.com/MustafaAmer-1/Flux/internal/database"
)

func (apiCfg *apiConfig) middlewareAuth(authHandler func(http.ResponseWriter, *http.Request, database.User)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusForbidden, fmt.Sprintf("Auth error: %v", err))
			return
		}
		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Auth error: no user with this apikey")
			return
		}

		authHandler(w, r, user)
	}
}
