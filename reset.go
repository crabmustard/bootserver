package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerResetUsers(w http.ResponseWriter, r *http.Request) {
	log.Print(cfg.platform + "\n")
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment"))
		return
	}
	cfg.fileServerHits.Store(0)
	cfg.db.Reset(r.Context())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and database reset to inital state."))
}
