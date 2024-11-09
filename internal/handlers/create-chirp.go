package handlers

import (
	"chirpy/internal/config"
	"chirpy/internal/database"
	"chirpy/internal/types"
	"chirpy/internal/utils"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func CreateChirpHandler(cfg *config.ApiConfig) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		type params struct {
			Body string `json:"body"`
		}

		data := params{}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&data); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
			return
		}

		userId, ok := utils.GetUserIDFromContext(r)

		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "Your are not authorized to see this page", nil)
			return
		}

		if len(data.Body) > 140 {
			utils.RespondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
			return
		}

		cleanedMessage := utils.ReplaceBadWords(data.Body, []string{"kerfuffle", "sharbert", "fornax"})

		newChirp := database.CreateChirpParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Body:      cleanedMessage,
			UserID:    userId,
		}

		dbChirp, err := cfg.Db.CreateChirp(r.Context(), newChirp)

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't write to db: %s", err)
			return
		}

		chirp := types.Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserId:    dbChirp.UserID,
		}
		utils.RespondWithJSON(w, http.StatusCreated, chirp)

	}
}
