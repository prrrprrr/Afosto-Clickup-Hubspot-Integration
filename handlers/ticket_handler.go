package handlers

import (
	"Afosto-Clickup-Hubspot-Integration/tickets"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

type TicketHandler struct {
	Service tickets.TicketService
}

func UpdateTicket(s TicketHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			sendError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
			return
		}
		ctx := r.Context()
		//parse request
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			sendError(w, http.StatusBadRequest, errors.New("error reading body"))
			return
		}
		defer func() { _ = r.Body.Close() }()
		var payload tickets.UpdateTicketWebhookRequest
		err = json.Unmarshal(body, &payload)
		if err != nil {
			log.Printf("error unmarshalling updateTicketRequest: %v", err)
			sendError(w, http.StatusBadRequest, errors.New("error unmarshalling updateTicketRequest"))
			return
		}
		log.Printf("updating ticketstatus for %s", payload.Payload.ID)
		//call service to update task
		err = s.Service.ProcessUpdateTicketStatus(ctx, payload)
		if err != nil {
			log.Printf("error updating ticket status for task #%s, with error: %s", payload.Payload.ID, err)
			sendError(w, http.StatusInternalServerError, errors.New("error updating ticket status"))
			return
		}

		sendJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"message": "ticket updated",
		})
	}
}
