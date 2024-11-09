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

	"github.com/google/uuid"
)

type params struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUserHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := params{}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&params); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
			return
		}

		password, err := auth.HashPassword(params.Password)

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		}

		newUser := database.CreateUserParams{
			Email:          params.Email,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			ID:             uuid.New(),
			HashedPassword: password,
		}

		dbUser, err := cfg.Db.CreateUser(r.Context(), newUser)

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error writing to db $s", err)
			return
		}

		user := types.User{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
		}

		utils.RespondWithJSON(w, http.StatusCreated, user)
	}
}
