package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"email"`
		UserID uuid.UUID `json:"user_id"`
	}
	type response struct {
		Chirp
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt decode parameters", err)
		return
	}
	dbChirp, err := cfg.db.CreateChirp(r.Context(), params.body, params.userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error making user", err)
		return
	}

	respondWithJson(w, http.StatusCreated, response{
		Chirp: Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		},
	})

}
