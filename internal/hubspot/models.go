package hubspot

import "time"

const (
	HubSpotGetTicketThreadsEndpoint = "https://api.hubapi.com/conversations/v3/conversations/threads?associatedTicketId=%s"
	HubSpotCRMEndPoint              = "https://api.hubapi.com/crm/v3/objects/0-5/%s" // Used to manipulate HubSpot Tickets, still requires hs Ticket FileID
	HubSpotPostToThreadEndpoint     = "https://api.hubapi.com/conversations/v3/conversations/threads/%s/messages"
	GetFileEndpoint                 = "https://api.hubapi.com/files/v3/files/%s"
	PostFileEndpoint                = "https://api.hubapi.com/files/v3/files"
	folderID                        = "390279679201"
	GetEngagementEndpoint           = "https://api.hubapi.com/engagements/v1/engagements/%s"
	HubSpotSignedUrlEndpoint        = "https://api-eu1.hubspot.com/files/v3/files/%s/signed-url"
)

type PipeLineResponse struct {
	Results []struct {
		Label        string `json:"label"`
		DisplayOrder int    `json:"displayOrder"`
		ID           string `json:"id"`
		Stages       []struct {
			Label        string `json:"label"`
			DisplayOrder int    `json:"displayOrder"`
			Metadata     struct {
				TicketState string `json:"ticketState"`
				IsClosed    string `json:"isClosed"`
			} `json:"metadata"`
			ID               string    `json:"id"`
			CreatedAt        time.Time `json:"createdAt"`
			UpdatedAt        time.Time `json:"updatedAt"`
			WritePermissions string    `json:"writePermissions"`
			Archived         bool      `json:"archived"`
		} `json:"stages"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
		Archived  bool      `json:"archived"`
	} `json:"results"`
}
