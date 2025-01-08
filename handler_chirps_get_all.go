package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerChirpsGetAll(w http.ResponseWriter, r *http.Request) {
	type response struct {
		allChirps []Chirp
	}

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
	chirpResponse := response{
		allChirps: jsonChirps,
	}
	log.Print(chirpResponse)
	respondWithJson(w, http.StatusOK, chirpResponse)

}
