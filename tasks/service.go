package tasks

import (
	"context"
	"fmt"
	"Afosto-Clickup-Hubspot-Integration/attachments"
	"Afosto-Clickup-Hubspot-Integration/Internal/clickup"
	"Afosto-Clickup-Hubspot-Integration/Internal/hubspot"
	"log"
)

type TaskService interface {
	ProcessCreateTask(ctx context.Context, payload CreateTaskPayload) (err error)
	ProcessUpdateTask(ctx context.Context, payload UpdateStatusPayload) (err error)
	leaveInitializationComment(ctx context.Context, clickupTaskID string, hubspotTicketID string) (response *clickup.LeaveCommentResponse, err error)
}
type Settings struct {
	ClickupClient              clickup.ClickUpClient
	HubspotClient              hubspot.HubSpotClient
	AttachmentService          attachments.AttachmentService
	HubSpotTicketCustomFieldID string
	CustomerCustomFieldID      string
	FallbackCustomerID         string
	HubSpotPortalID            string
	EmailContactID             string
}
type taskService struct {
	ClickupClient              clickup.ClickUpClient
	HubspotClient              hubspot.HubSpotClient
	AttachmentService          attachments.AttachmentService
	hubSpotTicketCustomFieldID string
	customerCustomFieldID      string
	fallbackCustomerID         string
	hubSpotPortalID            string
	emailContactID             string
}

//NewService is the factory for taskService
func NewService(
	settings Settings,
) TaskService {
	return &taskService{
		ClickupClient:              settings.ClickupClient,
		HubspotClient:              settings.HubspotClient,
		AttachmentService:          settings.AttachmentService,
		hubSpotTicketCustomFieldID: settings.HubSpotTicketCustomFieldID,
		customerCustomFieldID:      settings.CustomerCustomFieldID,
		fallbackCustomerID:         settings.FallbackCustomerID,
		hubSpotPortalID:            settings.HubSpotPortalID,
		emailContactID:             settings.EmailContactID,
	}
}

/*
ProcessUpdateTask handles the requests to update ClickUp tasks.
This is used when a change is made in HubSpot and this change needs to be reflected in ClickUp.
For now, it just changes status, but it can be expanded
*/
func (s taskService) ProcessUpdateTask(ctx context.Context, request UpdateStatusPayload) error {
	log.Printf("updating task status of task: #%s to: %s:", request.TaskID, request.Status)
	//compose body
	reqBody := clickup.UpdateTaskPayload{Status: request.Status}
	//do request
	err := s.ClickupClient.UpdateTask(ctx, request.TaskID, reqBody)
	if err != nil {
		return fmt.Errorf("error updating task status %w", err)
	}
	//return
	log.Printf("done updating status")
	return nil
}

/*
	ProcessCreateTask handles task creation for every HubSpot ticket that is created, filling it with relevant data after creation.
	It handles the inline and attached attachments so they're available in the Clickup Task.
	When it's done it leaves a comment with the link to the corresponding Hubspot ticket so navigation is easy.
*/
func (s taskService) ProcessCreateTask(ctx context.Context, payload CreateTaskPayload) error {

	//attempt to figure out the customer-id from clickup so we can sort the task under that customer
	customerID, err := s.ClickupClient.GetCustomerID(ctx, payload.Customer)
	if err != nil {
		log.Printf("error fetching matching id for customer: %v", err)
		customerID = s.fallbackCustomerID
	}

	//fill in custom field data for clickup request
	custom := []clickup.ClickupCustomField{
		{CustomFieldID: s.customerCustomFieldID, CustomFieldValue: customerID},
		{CustomFieldID: s.hubSpotTicketCustomFieldID, CustomFieldValue: payload.TicketID},
		{CustomFieldID: s.emailContactID, CustomFieldValue: payload.EmailFrom},
	}

	//fil request with data
	reqBody := clickup.CreateTaskRequest{
		Description:  fmt.Sprintf("%s\n%s", payload.EmailFrom, payload.Description),
		TicketName:   payload.TicketName,
		CustomFields: custom,
	}

	//call create_task endpoint
	resp, err := s.ClickupClient.CreateTask(ctx, reqBody)
	if err != nil {
		return fmt.Errorf("error calling create_task API: %w", err)
	}

	//handle attachments
	log.Printf("handling attachments for: %s", resp.TaskId)
	engagement, err := s.HubspotClient.GetEngagement(ctx, payload.EmailID)
	if err != nil {
		log.Printf("error pulling engagement for: %s with the error: %v", resp.TaskId, err)
	}
	err = s.AttachmentService.ProcessAttachments(ctx, resp.TaskId, *engagement)
	if err != nil {
		log.Printf("error handling attachments for: %s with the error: %v", resp.TaskId, err)
	}

	//leave comment in ClickUp task containing the HubSpot link
	//ClickUp will link the task to the ticket this way
	//we don't care about the response
	log.Printf("leaving comment on: %s", resp.TaskId)
	_, err = s.leaveInitializationComment(ctx, resp.TaskId, payload.TicketID)
	if err != nil {
		log.Printf("leaving comment for %s went wrong", resp.TaskId)
	}

	//now we send task id and url to the hubspot ticket so both objects have their counterpart's id
	//fil request with data
	log.Printf("sending info back to hubspot")
	req := hubspot.UpdateTicketRequest{
		Properties: hubspot.TicketProperties{
			ClickUpTaskid:  resp.TaskId,
			ClickUpTaskURL: fmt.Sprintf(clickUpTaskUrl, resp.TaskId),
		},
	}

	//we do not care about the response
	//if there is no error the comment was posted successfully, and we don't need to check anything else
	_, err = s.HubspotClient.UpdateTicket(ctx, payload.TicketID, req)
	if err != nil {
		return fmt.Errorf("error updating HubSpot ticket: %w", err)
	}
	return nil
}

/*
This function is used to leave a comment under a newly created ClickUp task which contains the hubspot ticket link. This creates a link between clickup and hubspot
*/
func (s taskService) leaveInitializationComment(ctx context.Context, clickUpTaskID string, hubSpotTicketID string) (*clickup.LeaveCommentResponse, error) {

	//compose body
	hubSpotLink := fmt.Sprintf(hubSpotTicketLink, s.hubSpotPortalID, hubSpotTicketID)
	reqBody := clickup.LeaveCommentRequest{CommentText: hubSpotLink, NotifyAll: false}

	//send request
	resp, err := s.ClickupClient.LeaveTaskComment(ctx, clickUpTaskID, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error leaving InitializationComment: %w", err)
	}

	//return the response
	return resp, nil
}

///*
//MoveToCurrentSprint
//*/
//func (s taskService) MoveToCurrentSprint(ctx context.Context, taskid string) error {
//	//get current sprint listID
//call endpoint to move task to that list id
//finish
//}
