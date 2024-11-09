package config

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"chirpy/internal/types"
	"chirpy/internal/utils"
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/google/uuid"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	Db             *database.Queries
	Platform       string
	JwtSecret      string
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) GetMetrics(w http.ResponseWriter, r *http.Request) {
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

func (cfg *ApiConfig) Reset(w http.ResponseWriter, r *http.Request) {
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

func (cfg *ApiConfig) AuthorizationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer, err := auth.GetBearerToken(r.Header)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Your are not authorized to see this page", err)
			return
		}
		userID, err := auth.ValidateJWT(bearer, cfg.JwtSecret)

		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Your are not authorized to see this page", err)
			return
		}

		ctx := context.WithValue(r.Context(), types.UserIDKey, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) GetUserIdFromToken(r *http.Request, tokenString string) (uuid.UUID, error) {
	userId, err := cfg.Db.GetUserByRefreshToken(r.Context(), tokenString)

	if err != nil {
		return uuid.Nil, err
	}

	return userId, errors.New("userId not found in token")
}
