package handlers

import (
	"Afosto-Clickup-Hubspot-Integration/tasks"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

type TaskHandler struct {
	Service tasks.TaskService
}

func CreateTask(s TaskHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//only allow Post method
		if r.Method != http.MethodPost {
			sendError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
			return
		}

		//read the body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("error unmarshalling payload: %v", err)
			sendError(w, http.StatusBadRequest, errors.New("error reading body"))
			return
		}
		defer func() { _ = r.Body.Close() }()

		//unmarshall request data
		var rootPayload tasks.CreateTaskWebhook
		if err = json.Unmarshal(body, &rootPayload); err != nil {
			log.Printf("error unmarshalling payload: %v", err)
			sendError(w, http.StatusBadRequest, errors.New("error unmarshalling body"))
			return
		}

		//prepare payload
		var payload tasks.CreateTaskPayload
		payload = rootPayload.Fields
		ctx := r.Context()

		//process request
		err = s.Service.ProcessCreateTask(ctx, payload)
		if err != nil {
			log.Printf("error creating task: %v", err)
			sendError(w, http.StatusInternalServerError, errors.New("error creating task"))
			return
		}
		//return response
		sendJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"message": "ticket created",
		})
	}
}

func UpdateTask(s TaskHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//only allow Post method
		if r.Method != http.MethodPost {
			sendError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
			return
		}
		//read the body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			sendError(w, http.StatusBadRequest, errors.New("Error reading body"))
			return
		}
		defer func() { _ = r.Body.Close() }()

		//fill payload struct with payload data
		var rootPayload tasks.UpdateTaskWebhook
		if err = json.Unmarshal(body, &rootPayload); err != nil {
			log.Printf("Error parsing body: %w", body)
			sendError(w, http.StatusBadRequest, errors.New("invalid JSON"))
			return
		}
		var payload tasks.UpdateStatusPayload
		payload = rootPayload.Fields
		ctx := r.Context()

		err = s.Service.ProcessUpdateTask(ctx, payload)
		if err != nil {
			log.Printf("Error processing task update : %v", err)
			sendError(w, http.StatusBadRequest, errors.New("error processing task update"))
			return
		}

		sendJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"message": "task updated",
		})
	}
}
