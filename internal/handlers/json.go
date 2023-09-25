package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		logrus.Warn("Responding with 5XX error: ", msg)
	}
	type errResponce struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errResponce{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal JSON responce: ", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
