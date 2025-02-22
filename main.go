package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"github.com/joho/godotenv"
	"os"
	"time"
	"database/sql"
	"github.com/V-Gira/chirpy/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits 	atomic.Int32
	db      		*database.Queries
	platform 		string
	secret			string
	timeout		time.Duration
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	JWTsecret := os.Getenv("JWT_SECRET")
	if JWTsecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	queries := database.New(db)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db : queries,
		platform: platform,
		secret: JWTsecret,
		timeout: time.Hour,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/login", apiCfg.handlerUsersLogin)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUsersPut)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsGetAll)
	mux.HandleFunc("GET /api/chirps/{chirpsID}", apiCfg.handlerChirpsGetOne)
	mux.HandleFunc("DELETE /api/chirps/{chirpsID}", apiCfg.handlerChirpsDelete)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRefreshRevoke)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}