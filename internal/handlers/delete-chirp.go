package handlers

import (
	"chirpy/internal/config"
	"chirpy/internal/utils"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func DeleteChirp(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirpId, err := uuid.Parse(r.PathValue("id"))

		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Couldn't parse param id: %s", err)
			return
		}

		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Not authorized", nil)
			return
		}

		chirp, err := cfg.Db.GetChirp(r.Context(), chirpId)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				utils.RespondWithError(w, http.StatusNotFound, "chirp not found", nil)
				return
			}
			utils.RespondWithError(w, http.StatusInternalServerError, "Error get chirp: %s", err)
			return
		}

		userId, ok := utils.GetUserIDFromContext(r)

		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "Your are not authorized to see this page", nil)
			return
		}

		fmt.Printf("user id %s", userId)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Not authorized", nil)
			return
		}

		if chirp.UserID != userId {
			utils.RespondWithError(w, http.StatusForbidden, "Not authorized", nil)
			return
		}

		err = cfg.Db.DeleteChirp(r.Context(), chirpId)

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp: %s", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
