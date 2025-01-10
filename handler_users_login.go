package main

import (
	"encoding/json"
	"net/http"

	"github.com/crabmustard/bootserver/internal/auth"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt decode parameters", err)
		return
	}
	dbUser, err := cfg.db.GetPasswordHash(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error making user", err)
		return
	}
	ok := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if ok != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	respondWithJson(w, http.StatusOK, User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	})

}
