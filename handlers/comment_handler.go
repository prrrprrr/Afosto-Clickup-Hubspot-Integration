package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"Afosto-Clickup-Hubspot-Integration/comments"
	"log"
	"net/http"
	"os"
	"strings"
)

type CommentHandler struct {
	Service comments.CommentService
}

/*
Handles logic for inbound emails
*/

func ReceiveEmail(s CommentHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req comments.RecieveEmailWebhook
		ctx := r.Context()
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("error decoding body: %v", err)
			sendError(w, http.StatusBadRequest, errors.New("bad request"))
			return
		}
		payload := comments.SendEmailToClickup{
			ClickUpTaskID: req.Fields.ClickupTaskID,
			EngagementID:  req.Fields.EmailID,
			TicketStatus:  strings.ToLower(req.Fields.TicketStatus),
			TicketId:      req.Fields.TicketId,
		}
		err := s.Service.ProcessEmailFromCustomer(ctx, payload)
		if err != nil {
			log.Printf("could not process incoming email: %s, with error: %w", payload.EngagementID, err)
			sendError(w, http.StatusInternalServerError, errors.New("could not process email"))
			return
		}
		sendJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"message": "comment webhook processed",
			"data":    req,
		})
	}
}

/*
SendEmail handles outbound email logic
*/
func SendEmail(s CommentHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req comments.SendEmailWebhookRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("error decoding request: %v", err)
			sendJSON(w, http.StatusBadRequest, errors.New("bad request"))
			return
		}
		var ticketID string
		//get the ticketID using the clickup custom field
		for _, v := range req.Payload.Fields {
			if v.FieldID == os.Getenv("CLICKUP_CUSTOM_FIELD_ID_HUBSPOT_TICKET_ID") {
				ticketID = fmt.Sprintf("%s", v.Value)
			}
		}

		payload := comments.SendEmailToCustomer{
			TaskID:   req.Payload.ID,
			TicketID: ticketID,
		}
		//process request
		err := s.Service.ProcessEmailToCustomer(ctx, payload)
		if err != nil {
			log.Printf("error processing email to customer: %v", err)
			sendError(w, http.StatusInternalServerError, errors.New("something went wrong"))
			return
		}

		sendJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"message": "comment webhook processed",
			"data":    req,
		})
	}
}
