package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func sendJSON(w http.ResponseWriter, status int, payload any) {
	//respond to webhook sender
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("sendJSON: encode error: %v", err)
	}

}

func sendError(w http.ResponseWriter, status int, err error) {
	sendJSON(w, status, struct {
		Error string `json:"error"`
	}{err.Error()})
}
