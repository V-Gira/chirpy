package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"github.com/joho/godotenv"
	"os"
	"database/sql"
	"github.com/V-Gira/chirpy/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db      *database.Queries
}

func main() {
	godotenv.Load()

	const filepathRoot = "."
	const port = "8080"
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	queries := database.New(db)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db : queries,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}