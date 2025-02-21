package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"github.com/google/uuid"
	"github.com/V-Gira/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type response struct {
		Chirp
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

	chirpsParams := database.CreateChirpParams{
		Body: params.Body,
		UserID: params.UserID,
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), chirpsParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Chirp: Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		},
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
