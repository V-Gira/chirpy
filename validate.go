package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	params.Body = replaceProfanity(params.Body)

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: params.Body,
	})
}

func replaceProfanity(s string) string {
	profanities := []string{"kerfuffle", "sharbert", "fornax"}
	wordSplit := strings.Split(s, " ")
	for i, word := range wordSplit {
		for _, profanity := range profanities {
			if strings.ToLower(word) == strings.ToLower(profanity) {
				wordSplit[i] = "****"
			}
		}
	}
	return strings.Join(wordSplit, " ")
}
