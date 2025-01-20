package main

import (
	"encoding/json"
	"net/http"

	"github.com/crabmustard/bootserver/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpyRed(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		}
	}
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no api key in header", err)
		return
	}
	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "invalid api key", nil)
		return

	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt decode parameters", err)
		return
	}
	if params.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "not a proper event", nil)
		return
	}
	upgradeErr := cfg.db.EnableChirpyRed(r.Context(), params.Data.UserID)
	if upgradeErr != nil {
		respondWithError(w, http.StatusNotFound, "user not found", upgradeErr)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
