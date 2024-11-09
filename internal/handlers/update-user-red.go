package handlers

import (
	"chirpy/internal/auth"
	"chirpy/internal/config"
	"chirpy/internal/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type EventData struct {
	UserId string `json:"user_id"`
}

type UpdateUserRedParams struct {
	Event string    `json:"event"`
	Data  EventData `json:"data"`
}

func UpdateUserRed(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey, err := auth.GetAPIKey(r.Header)

		if err != nil || apiKey != os.Getenv("POLKA_KEY") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		updateUserParams := UpdateUserRedParams{}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&updateUserParams); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
			return
		}

		if updateUserParams.Event != "user.upgraded" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		userId, err := uuid.Parse(updateUserParams.Data.UserId)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = cfg.Db.UpdateUserRed(r.Context(), userId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
