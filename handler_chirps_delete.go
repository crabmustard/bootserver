package main

import (
	"net/http"

	"github.com/crabmustard/bootserver/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {

	jwToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to find bearer token", err)
		return
	}
	user, err := auth.ValidateJWT(jwToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "unable to find bearer token", err)
		return
	}

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error parsing chirpid", err)
		return
	}
	chirp, err := cfg.db.GetChirpById(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "unable to find chirp in db", err)
		return
	}
	if chirp.UserID != user {
		respondWithError(w, http.StatusForbidden, "userid dont match", err)
		return
	}
	deleted := cfg.db.DeleteChirpById(r.Context(), chirpID)
	if deleted != nil {
		respondWithError(w, http.StatusNotFound, "unable to delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
