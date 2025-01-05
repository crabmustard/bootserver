package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type validResp struct {
		Valid bool   `json:"valid"`
		Chirp string `json:"cleaned_body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt decode params")
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	cleaned := cleanChirp(params.Body)

	respondWithJson(w, 200, validResp{
		Valid: true,
		Chirp: cleaned,
	})

}

func cleanChirp(chirpBody string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	splitChirp := strings.Split(chirpBody, " ")
	for i, word := range splitChirp {
		lowWord := strings.ToLower(word)
		if slices.Contains(badWords, lowWord) {
			splitChirp[i] = "****"
		}
	}
	cleaned := strings.Join(splitChirp, " ")
	return cleaned
}
