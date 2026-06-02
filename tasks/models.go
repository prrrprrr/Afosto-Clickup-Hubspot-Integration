package tasks

const (
	hubSpotTicketLink = "https://app-eu1.hubspot.com/contacts/%s/record/0-5/%s" // HubspotTicket endpoint, needs hs portal FileID and hs ticket FileID
	clickUpTaskUrl    = "https://app.clickup.com/t/%s"                          // URL of a ClickUp task, needs ClickUp Task FileID
)

type UpdateStatusPayload struct {
	TaskID          string `json:"clickup_id"`
	Status          string `json:"ticket_status,omitempty"`
	MarkdownContent string `json:"markdown_content,omitempty"`
	Description     string `json:"description,omitempty"`
}

// webhookPayload represents the incoming webhook data.
type CreateTaskPayload struct {
	TicketID    string `json:"ticket_id"`
	TicketName  string `json:"ticket_name"`
	Customer    string `json:"clickup_customer,omitempty"`
	Description string `json:"ticket_description"`
	EmailBody   string `json:"email_body"`
	EmailFrom   string `json:"email_from"`
	EmailID     string `json:"email_id,omitempty"`
}
type CreateTaskWebhook struct {
	Fields CreateTaskPayload `json:"fields"`
}
type UpdateTaskWebhook struct {
	Fields UpdateStatusPayload `json:"fields"`
}

var PipelineToStatus = map[string]string{
	"1":          "inbox",
	"2":          "wachten op reactie",
	"3":          "terugkoppelen",
	"4":          "complete",
	"5302752480": "to do",
	"5302752481": "in progress",
}
