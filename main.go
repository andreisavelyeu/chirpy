package main

import (
	"chirpy/internal/config"
	"chirpy/internal/database"
	"chirpy/internal/handlers"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var port = "8080"

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("TOKEN_SECRET")
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Printf("database connection failed %s", err)
		return
	}

	dbQueries := database.New(db)

	config := &config.ApiConfig{
		Db:        dbQueries,
		Platform:  platform,
		JwtSecret: jwtSecret,
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	mux.Handle("/app/", http.StripPrefix("/app/", config.MiddlewareMetricsInc(http.FileServer(http.Dir("./")))))
	mux.Handle("/app/assets/", http.StripPrefix("/app/assets/", config.MiddlewareMetricsInc(http.FileServer(http.Dir("./assets")))))

	mux.HandleFunc("GET /api/healthz", handlers.HealthzHandler)
	mux.HandleFunc("GET /admin/metrics", config.GetMetricsHandler)
	mux.HandleFunc("POST /admin/reset", config.ResetHandler)
	mux.HandleFunc("POST /api/chirps", config.AuthorizationMiddleware(handlers.CreateChirpHandler(config)))
	mux.HandleFunc("POST /api/users", handlers.CreateUserHandler(config))
	mux.HandleFunc("GET /api/chirps", handlers.GetChirps(config))
	mux.HandleFunc("GET /api/chirps/{id}", handlers.GetChirp(config))
	mux.HandleFunc("POST /api/login", handlers.Login(config))
	mux.HandleFunc("POST /api/refresh", handlers.Refresh(config))
	mux.HandleFunc("POST /api/revoke", handlers.Revoke(config))

	fmt.Printf("Listening on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}
