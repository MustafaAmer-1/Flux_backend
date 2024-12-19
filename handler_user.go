package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MustafaAmer-1/Flux/internal/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	params := struct {
		Name   string `json:"name"`
		Passwd string `json:"password"`
	}{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Passwd), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't create a user: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Passwd:    hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't create a user: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	params := struct {
		Name   string `json:"name"`
		Passwd string `json:"password"`
	}{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.GetUserByName(r.Context(), params.Name)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't login a user: Invalid Name Or Password!")
		return
	}
	err = bcrypt.CompareHashAndPassword(user.Passwd, []byte(params.Passwd))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't login a user: Invalid Name Or Password!")
		return
	}
	respondWithJSON(w, http.StatusCreated, databaseUserToUser(user))
}
