package handlers

import (
	"chirpy/internal/auth"
	"chirpy/internal/config"
	"chirpy/internal/database"
	"chirpy/internal/utils"
	"net/http"
	"time"
)

func Revoke(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)

		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Not authorized", nil)
			return
		}

		params := database.RevokeRefreshTokenParams{
			UpdatedAt: time.Now(),
			Token:     token,
		}

		if err := cfg.Db.RevokeRefreshToken(r.Context(), params); err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Not authorized", nil)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
