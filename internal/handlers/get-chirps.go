package handlers

import (
	"chirpy/internal/config"
	"chirpy/internal/types"
	"chirpy/internal/utils"
	"net/http"
)

func GetChirps(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbChirps, err := cfg.Db.GetChirps(r.Context())

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't get from db: %s", err)
			return
		}

		chirps := make([]types.Chirp, len(dbChirps))

		for k, v := range dbChirps {
			chirps[k] = types.Chirp{
				ID:        v.ID,
				CreatedAt: v.CreatedAt,
				UpdatedAt: v.UpdatedAt,
				Body:      v.Body,
				UserId:    v.UserID,
			}
		}

		utils.RespondWithJSON(w, http.StatusOK, chirps)
	}
}
