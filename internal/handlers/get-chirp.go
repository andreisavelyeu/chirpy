package handlers

import (
	"chirpy/internal/config"
	"chirpy/internal/types"
	"chirpy/internal/utils"
	"net/http"

	"github.com/google/uuid"
)

func GetChirp(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := uuid.Parse(r.PathValue("id"))

		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Couldn't parse param id: %s", err)
			return
		}

		dbChirp, err := cfg.Db.GetChirp(r.Context(), userId)

		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Couldn't get from db: %s", err)
			return
		}

		chirp := types.Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserId:    dbChirp.UserID,
		}

		utils.RespondWithJSON(w, http.StatusOK, chirp)
	}
}
