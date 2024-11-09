package handlers

import (
	"chirpy/internal/auth"
	"chirpy/internal/config"
	"chirpy/internal/database"
	"chirpy/internal/types"
	"chirpy/internal/utils"
	"encoding/json"
	"net/http"
	"time"
)

func Login(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type params struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		data := params{}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&data); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
			return
		}

		dbUser, err := cfg.Db.GetUser(r.Context(), data.Email)

		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
			return
		}

		if err := auth.CheckPasswordHash(data.Password, dbUser.HashedPassword); err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
			return
		}

		tokenString, err := auth.MakeJWT(dbUser.ID, cfg.JwtSecret, time.Second*time.Duration(3600))

		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
			return
		}

		refreshToken, err := auth.MakeRefreshToken()

		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
			return
		}

		tokenParams := database.CreateRefreshTokenParams{
			Token:     refreshToken,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    dbUser.ID,
			ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		}

		dbRefreshToken, err := cfg.Db.CreateRefreshToken(r.Context(), tokenParams)

		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
			return
		}

		user := types.User{
			ID:           dbUser.ID,
			CreatedAt:    dbUser.CreatedAt,
			UpdatedAt:    dbUser.UpdatedAt,
			Email:        dbUser.Email,
			Token:        tokenString,
			RefreshToken: dbRefreshToken.Token,
			IsChirpyRed:  dbUser.IsChirpyRed,
		}

		utils.RespondWithJSON(w, http.StatusOK, user)
	}
}
