package config

import (
	"chirpy/internal/database"
	"chirpy/internal/utils"
	"fmt"
	"net/http"
	"sync/atomic"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	Db             *database.Queries
	Platform       string
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) GetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	hits := cfg.FileserverHits.Load()
	responseText := fmt.Sprintf(`
		<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>
	`, hits)
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(200)
	w.Write([]byte(responseText))
}

func (cfg *ApiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		utils.RespondWithError(w, http.StatusForbidden, "Forbidden", nil)
		return
	}

	err := cfg.Db.DeleteAllUsers(r.Context())

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Delete users error: %s\n", err)
		return
	}

	cfg.FileserverHits.Store(0)
	hits := cfg.FileserverHits.Load()
	responseText := fmt.Sprintf("Hits: %v\n", hits)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(responseText))
}