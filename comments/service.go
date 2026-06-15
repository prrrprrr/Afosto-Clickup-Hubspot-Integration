package comments

/*
This file contains the logic for processing e-mails from, and to the customer.
*/
import (
	"Afosto-Clickup-Hubspot-Integration/attachments"
	"Afosto-Clickup-Hubspot-Integration/internal/clickup"
	"Afosto-Clickup-Hubspot-Integration/internal/hubspot"
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
)

type CommentService interface {
	ProcessEmailFromCustomer(ctx context.Context, payload SendEmailToClickup) error
	ProcessEmailToCustomer(ctx context.Context, payload SendEmailToCustomer) error
}
type Settings struct {
	ClickupClient     clickup.ClickUpClient
	HubspotClient     hubspot.HubSpotClient
	AttachmentService attachments.AttachmentService
	ActorID           string
}
type commentService struct {
	ClickupClient     clickup.ClickUpClient
	HubspotClient     hubspot.HubSpotClient
	AttachmentService attachments.AttachmentService
	ActorID           string
}

//NewService is the factory for taskService
func NewService(
	settings Settings,
) CommentService {
	return commentService{
		ClickupClient:     settings.ClickupClient,
		HubspotClient:     settings.HubspotClient,
		AttachmentService: settings.AttachmentService,
		ActorID:           settings.ActorID,
	}
}

/*
ProcessEmailToCustomer handles newComment event in clickup
gets the newest comment, checks if it needs to be sent out, and does so if required
*/
func (s commentService) ProcessEmailToCustomer(ctx context.Context, payload SendEmailToCustomer) error {

	//collect and parse the last comment from the clickup task given in the payload
	log.Printf("sending email for task: %s", payload.TaskID)
	parsedComment, err := s.ClickupClient.GetLastComment(ctx, payload.TaskID)
	if err != nil {
		return fmt.Errorf("error pulling last comment from ClickUp task: %w", err)
	}

	if strings.Contains(parsedComment.Comment, "<p>@klant</p>") {
		//parsedComment.Comment = strings.TrimSpace(strings.TrimPrefix(parsedComment.Comment, "<p>@klant</p>"))
		parsedComment.Comment = strings.Replace(parsedComment.Comment, "<p>@klant</p>", "", 1)
	} else {
		log.Printf("this comment was not for the customer")
		return nil
	}

	//pull task for target email and ticket_id
	resp, err := s.ClickupClient.GetTask(ctx, payload.TaskID)
	if err != nil {
		return fmt.Errorf("error pulling task from ClickUp: %w", err)
	}

	//pull the thread from hubspot api using the hubspot ticket id
	ticketThreads, err := s.HubspotClient.GetTicketThreads(ctx, payload.TicketID)
	if err != nil {
		return fmt.Errorf("error getting ticket threads from HubSpot: %w", err)
	}

	//we collect the contact email and the hubspot thread ID from the task
	var targetEmail string
	for _, v := range resp.CustomFields {
		if v.ID == clickup.EmailContactCustomfield {
			if strings.TrimSpace(string(v.Value)) != "" { //if the target email is filled in we use it to send the email
				targetEmail = string(v.Value)
			} else { //else we just take the creator email
				targetEmail = resp.Creator.Email
			}
		}
	}
	actorID, ok := EmailToActorId[parsedComment.User]
	if !ok {
		//if no correct id was found we have this fallback
		actorID = "23951706"
	}

	var recipients = []hubspot.Recipients{}
	//compose recipients
	recipients = append(recipients, hubspot.Recipients{
		DeliveryIdentifier: hubspot.DeliveryIdentifier{
			Type:  "HS_EMAIL_ADDRESS",
			Value: targetEmail,
		},
		Name:           targetEmail,
		RecipientField: "TO",
	})

	thread := ticketThreads.Threads[0]
	req := hubspot.HubSpotSendEmailRequest{
		Type: "MESSAGE",
		Text: parsedComment.Comment,
		Attachments: []struct {
			FileID string `json:"fileId,omitempty"`
			Type   string `json:"type,omitempty"`
		}{},
		Recipients:       recipients,
		SenderActorID:    fmt.Sprintf("A-%s", actorID),
		ChannelID:        thread.OriginalChannelID,
		ChannelAccountID: thread.OriginalChannelAccountID,
		RichText:         parsedComment.Comment,
		Subject:          fmt.Sprintf("%s (#%s)", resp.Name, payload.TicketID),
	}
	//push the message to the correct thread
	_, err = s.HubspotClient.PostToTicketThread(ctx, thread.ID, &req)
	if err != nil {
		return fmt.Errorf("error posting message to ticket thread: %w", err)
	}
	//now we add a threaded comment to the parsed clickup comment to notify the developer his mail was sent out.
	request := clickup.LeaveCommentRequest{
		NotifyAll:   false,
		CommentText: "sent",
	}
	_, err = s.ClickupClient.LeaveThreadedComment(ctx, parsedComment.ID, request)
	if err != nil {
		return fmt.Errorf("error posting Threaded comment to clickup: %w", err)
	}
	return nil
}

/*
ProcessEmailFromCustomer pulls an engagement from a hubspot ticket and pushes the message to a clickup in the form of a task comment
if the ticket we recieved
*/
func (s commentService) ProcessEmailFromCustomer(ctx context.Context, payload SendEmailToClickup) error {
	//pull engagement from hubspot engagement API
	engagementResp, err := s.HubspotClient.GetEngagement(ctx, payload.EngagementID)
	if err != nil {
		return fmt.Errorf("error requesting engagement from HubSpot: %w", err)
	}

	//upload the attachments that came with the engagement
	err = s.AttachmentService.ProcessAttachments(ctx, payload.ClickUpTaskID, *engagementResp)
	if err != nil {
		return fmt.Errorf("error processing attachments : %w", err)
	}

	//create comment request to send to the task.
	request := clickup.LeaveCommentRequest{
		NotifyAll:   false,
		CommentText: s.stripEmailThread(engagementResp.Metadata.Text), // remove threaded history from metadata.Text so only the recent email is left
	}

	//send the comment to the task_id we got from the payload
	_, err = s.ClickupClient.LeaveTaskComment(ctx, payload.ClickUpTaskID, request)
	if err != nil {
		return fmt.Errorf("error leaving comment : %w", err)
	}

	//for certain status's we want to do some extra updating
	if payload.TicketStatus == clickup.CompleteStatus {

		//if the target ticket of te email is 'complete' we set both ticket and task to 'reopened'
		err = s.setObjectsToStatus(ctx, clickup.OpnieuwGeopendStatus, payload.ClickUpTaskID, payload.TicketId)
		if err != nil {
			return fmt.Errorf("error updating clickupTask : %w", err)
		}

	} else if payload.TicketStatus == clickup.WachtenOpReactieStatus {

		//if we were waiting on a reply we set the status to reply received
		err = s.setObjectsToStatus(ctx, clickup.ReactieOntvangenStatus, payload.ClickUpTaskID, payload.TicketId)
		if err != nil {
			return fmt.Errorf("error updating hubspotTicket : %w", err)
		}

	}

	return nil
}
func (s commentService) setObjectsToStatus(ctx context.Context, status string, taskid string, ticketid string) error {

	//compose clickup update request
	reqBody := clickup.UpdateTaskPayload{Status: status}

	//do request
	err := s.ClickupClient.UpdateTask(ctx, taskid, reqBody)
	if err != nil {
		return fmt.Errorf("error updating task status %w", err)
	}

	//compose hubspot update request
	pipelineStatus, err := s.HubspotClient.StatusToPipeLine(status)
	if err != nil {
		return fmt.Errorf("error matching status to pipeline %w", err)
	}
	properties := hubspot.TicketProperties{
		TicketStatus: pipelineStatus,
	}

	//do request
	hubspotReq := hubspot.UpdateTicketRequest{Properties: properties}
	_, err = s.HubspotClient.UpdateTicket(ctx, ticketid, hubspotReq)
	if err != nil {
		return fmt.Errorf("error updating ticket status %w", err)
	}
	return nil
}

/*
stripEmailThread stript de threadhistory van de emailtext
zodra een patroon word gevonden die lijkt op: "Op vr 8 mei 2026 om 13:17 schreef Afosto <support@afosto.com>:" word de email afgekapt.
*/
func (s commentService) stripEmailThread(body string) string {

	patterns := []string{
		`(?m)^Op .{5,100} schreef .+?:`,
		`(?m)^On .{5,100} wrote:`,
		`(?m)^Van: .+`,
		`(?m)^From: .+`,
		`(?m)^>`,
		`(?m)^-{2,}`,
	}
	earliest := len(body)
	for _, p := range patterns {
		re := regexp.MustCompile(p)
		loc := re.FindStringIndex(body)
		if loc != nil && loc[0] < earliest {
			earliest = loc[0]
		}
	}
	return strings.TrimRight(body[:earliest], "\r\n ")
}

type SendEmailToCustomer struct {
	TaskID   string
	TicketID string
}
type SendEmailToClickup struct {
	EngagementID  string `json:"engagement_id"`
	ClickUpTaskID string `json:"clickup_task_id"`
	TicketStatus  string `json:"ticket_status"`
	TicketId      string `json:"ticket_id"`
}
