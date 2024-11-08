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
			Body   string `json:"body"`
			UserId string `json:"user_id"`
		}

		data := params{}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&data); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
			return
		}

		if len(data.Body) > 140 {
			utils.RespondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
			return
		}

		cleanedMessage := utils.ReplaceBadWords(data.Body, []string{"kerfuffle", "sharbert", "fornax"})

		userId, err := uuid.Parse(data.UserId)

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't parse userId: %s", err)
			return
		}

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
