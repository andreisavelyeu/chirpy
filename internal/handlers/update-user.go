package handlers

import (
	"chirpy/internal/auth"
	"chirpy/internal/config"
	"chirpy/internal/database"
	"chirpy/internal/types"
	"chirpy/internal/utils"
	"encoding/json"
	"net/http"
)

type updateUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func UpdateUser(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := updateUserParams{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&params); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
			return
		}

		userId, ok := utils.GetUserIDFromContext(r)

		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "Your are not authorized to see this page", nil)
			return
		}

		newUser := database.UpdateUserParams{
			ID: userId,
		}

		if params.Password != "" {
			password, err := auth.HashPassword(params.Password)
			if err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
			}
			newUser.HashedPassword = password
		}

		if params.Email != "" {
			newUser.Email = params.Email
		}

		updatedUser, err := cfg.Db.UpdateUser(r.Context(), newUser)

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Update failed", err)
		}

		user := types.User{
			ID:          updatedUser.ID,
			CreatedAt:   updatedUser.CreatedAt,
			UpdatedAt:   updatedUser.UpdatedAt,
			Email:       updatedUser.Email,
			IsChirpyRed: updatedUser.IsChirpyRed,
		}

		utils.RespondWithJSON(w, http.StatusOK, user)
	}
}
