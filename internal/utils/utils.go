package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func RespondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}

	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	RespondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}

func ReplaceBadWords(word string, badWords []string) string {
	splittedWord := strings.Split(word, " ")

	for i, v := range splittedWord {
		for _, badWord := range badWords {
			if strings.ToLower(v) == badWord {
				splittedWord[i] = "****"
			}
		}
	}
	return strings.Join(splittedWord, " ")
}
