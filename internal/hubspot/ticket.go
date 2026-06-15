package hubspot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

/*
UpdateTicket updates HubSpot ticket properties.
Expects an UpdateTicketRequest which contains properties one wishes to change within a ticket.
*/
func (s hubSpotClient) UpdateTicket(ctx context.Context, ticketID string, request UpdateTicketRequest) (*UpdateTicketResponse, error) {
	//compose url
	endPoint := fmt.Sprintf(HubSpotCRMEndPoint, ticketID)

	//compose body
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request that updates a HubSpot ticket: %w", err)
	}

	//do request
	resp, err := s.call(ctx, http.MethodPatch, endPoint, jsonData)
	if err != nil {
		return nil, fmt.Errorf("error sending request that updates a HubSpot ticket: %w", err)
	}

	//parse response
	var result UpdateTicketResponse
	json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing response from engagement API: %w", err)
	}

	//return parsed response
	return &result, nil
}

/*
StatusToPipeLine is a helper function for UpdateTicket. The clickup API and the Hubspot API both return status's as strings while only clickup accepts them.
This function fills the hole by turning the string into the right id that hubspor requires
*/

func (s hubSpotClient) StatusToPipeLine(status string) (string, error) {
	stat := strings.ToLower(status)
	res, ok := statusToPipeline[stat]
	if !ok {
		return "", fmt.Errorf("error matching pipeline to status")
	}
	return res, nil
}

var statusToPipeline = map[string]string{
	"inbox":              "1",
	"wachten op reactie": "2",
	"terugkoppelen":      "3",
	"complete":           "4",
	"to do":              "5302752480",
	"in progress":        "5302752481",
	"opnieuw geopend":    "5316325601",
	"reactie ontvangen":  "5350206711",
}

type UpdateTicketRequest struct {
	Properties TicketProperties `json:"properties"`
}
type TicketProperties struct {
	ClickUpTaskid  string `json:"clickup_taskid,omitempty"`
	ClickUpTaskURL string `json:"clickup_task_url,omitempty"`
	TicketStatus   string `json:"hs_pipeline_stage,omitempty"`
}
type TicketProperty struct {
	Key   string
	Value string
}
type UpdateTicketResponse struct {
	ID               string           `json:"id,omitempty"`
	TicketProperties TicketProperties `json:"properties,omitempty"`
}
