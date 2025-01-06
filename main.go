package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/crabmustard/bootserver/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("error opening db")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("Platform must be set")
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileServerHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
	}
	mux := http.NewServeMux()
	mux.Handle("/app/",
		apiCfg.middlewareMetricsInc(http.StripPrefix("/app",
			http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)

	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerResetUsers)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerGetMetrics)

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

}
