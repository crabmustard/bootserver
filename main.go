package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

const filepathRoot = "."
const port = "8080"

func main() {
	apiCfg := apiConfig{}
	mux := http.NewServeMux()
	mux.Handle("/app/",
		apiCfg.middlewareMetricsInc(http.StripPrefix("/app",
			http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /healthz", handlerReadiness)
	mux.HandleFunc("GET /metrics", apiCfg.handlerGetMetrics)
	mux.HandleFunc("POST /reset", apiCfg.handlerResetMetrics)

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
