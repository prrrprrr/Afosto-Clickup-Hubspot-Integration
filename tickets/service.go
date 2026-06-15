package tickets

import (
	"context"
	"fmt"
	"Afosto-Clickup-Hubspot-Integration/internal/clickup"
	"Afosto-Clickup-Hubspot-Integration/internal/hubspot"
	"log"
)

type TicketService interface {
	ProcessUpdateTicketStatus(ctx context.Context, payload UpdateTicketWebhookRequest) error
}

type ticketService struct {
	hubspotClient hubspot.HubSpotClient
	clickupClient clickup.ClickUpClient
}

func NewService(hubspotClient hubspot.HubSpotClient, clickupClient clickup.ClickUpClient) ticketService {
	return ticketService{
		hubspotClient: hubspotClient,
		clickupClient: clickupClient,
	}
}
func (s ticketService) ProcessUpdateTicketStatus(ctx context.Context, payload UpdateTicketWebhookRequest) error {
	//pull task from clickup for more information since the webhook does not give us everything we need
	taskResp, err := s.clickupClient.GetTask(ctx, payload.Payload.ID)
	if err != nil {
		return fmt.Errorf("error pulling task from clickup API: %w", err)
	}

	//figure out the ticket ID and status
	var ticketID string
	for _, v := range taskResp.CustomFields {
		if v.ID == clickup.HubSpotTicketIdCustomField {
			ticketID = string(v.Value)
		}
	}
	//if the hubspot ticket id field was not filled in we cannot update the ticket
	if ticketID == "" {
		return fmt.Errorf("error finding hubspot ticket id in custom fields. Is it filled in? #%s", payload.Payload.ID)
	}

	//save task status
	status, err := s.hubspotClient.StatusToPipeLine(taskResp.Status.Status)
	if err != nil {
		return fmt.Errorf("error matching pipeline stage to a status")
	}

	//compose request
	properties := hubspot.TicketProperties{
		TicketStatus: status,
	}
	req := hubspot.UpdateTicketRequest{Properties: properties}

	//send request
	_, err = s.hubspotClient.UpdateTicket(ctx, ticketID, req)
	if err != nil {
		return fmt.Errorf("error updating ticketstatus in hubspot: %w", err)
	}
	log.Printf("ticketstatus updated")

	//return nil if there were no errors
	return nil
}
