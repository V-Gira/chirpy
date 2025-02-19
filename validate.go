package main

import (
	"net/http"
	"log"
	"encoding/json"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(`{"error": "Something went wrong"}`))
		return
	}

	type returnVals struct {
		valid bool `json:"valid"`
	}

	respBody := returnVals{
		valid: len(params.Body) <= 140,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error decoding parameters: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(`{"error": "Something went wrong"}`))
		w.Write(dat)
	}
	status := 400
	 if respBody.valid {
		status = 200 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(`{"valid": true}`))

		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			w.Write([]byte(`{"error": "Chirp is too long"}`))
		}
}