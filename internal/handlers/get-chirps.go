package handlers

import (
	"chirpy/internal/config"
	"chirpy/internal/types"
	"chirpy/internal/utils"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func GetChirps(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		author_id := r.URL.Query().Get("author_id")
		sortBy := r.URL.Query().Get("sort")

		dbChirps, err := cfg.Db.GetChirps(r.Context())

		authorID := uuid.Nil
		if author_id != "" {
			authorID, err = uuid.Parse(author_id)
			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
				return
			}
		}

		sortDirection := "asc"
		if sortBy == "desc" {
			sortDirection = "desc"
		}

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't get from db: %s", err)
			return
		}

		chirps := make([]types.Chirp, len(dbChirps))

		for k, v := range dbChirps {

			if authorID != uuid.Nil && v.UserID != authorID {
				continue
			}

			chirps[k] = types.Chirp{
				ID:        v.ID,
				CreatedAt: v.CreatedAt,
				UpdatedAt: v.UpdatedAt,
				Body:      v.Body,
				UserId:    v.UserID,
			}
		}

		sort.Slice(chirps, func(i, j int) bool {
			if sortDirection == "desc" {
				return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
			}
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})

		utils.RespondWithJSON(w, http.StatusOK, chirps)
	}
}
