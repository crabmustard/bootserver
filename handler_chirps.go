package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGetAll(w http.ResponseWriter, r *http.Request) {

	allChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error making user", err)
		return
	}

	jsonChirps := []Chirp{}
	for _, chirp := range allChirps {
		nextChirp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			UserID:    chirp.UserID,
			Body:      chirp.Body,
		}
		jsonChirps = append(jsonChirps, nextChirp)
	}

	respondWithJson(w, http.StatusOK, jsonChirps)

}

func (cfg *apiConfig) handlerChirpsGetID(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error parsing chirpid", err)
		return
	}
	theChirp, err := cfg.db.GetChirpById(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error retreiving chirp", err)
		return
	}

	respondWithJson(w, http.StatusOK, Chirp{
		ID:        theChirp.ID,
		CreatedAt: theChirp.CreatedAt,
		UpdatedAt: theChirp.UpdatedAt,
		UserID:    theChirp.UserID,
		Body:      theChirp.Body,
	})

}
