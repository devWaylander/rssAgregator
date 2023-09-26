package main

import (
	"context"
	"encoding/json"
	"net/http"
)

type Repo interface {
	NewUser(ctx context.Context, name string) (error, bool)
}

func (rC *repo) handlerUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Error parsing JSON: ")
		return
	}

	rC.NewUser(r.Context(), params.Name)
}
