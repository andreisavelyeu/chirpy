package handlers

import (
	"chirpy/internal/auth"
	"chirpy/internal/config"
	"chirpy/internal/utils"
	"net/http"
	"time"
)

type ResponseBody struct {
	Token string `json:"token"`
}

func Refresh(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)

		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Not authorized", nil)
			return
		}

		userId, err := cfg.GetUserIdFromToken(r, token)

		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Not authorized", nil)
			return
		}

		tokenString, err := auth.MakeJWT(userId, cfg.JwtSecret, time.Second*time.Duration(3600))

		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Not authorized", nil)
			return
		}

		responseBody := ResponseBody{
			Token: tokenString,
		}

		utils.RespondWithJSON(w, http.StatusOK, responseBody)
	}
}
