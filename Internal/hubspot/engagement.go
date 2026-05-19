package hubspot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

/*
GetEngagement retrieves a single HubSpot engagement based on the FileID provided
*/
func (s hubSpotClient) GetEngagement(ctx context.Context, engagementID string) (response *EngagementResponse, err error) {
	//prepare endpoint
	endPoint := fmt.Sprintf(GetEngagementEndpoint, engagementID)
	//do request
	resp, err := s.call(ctx, http.MethodGet, endPoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error calling engagement api: %w", err)
	}

	//parse response
	var result EngagementResponse
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing response from engagement API: %w", err)
	}

	//return parsed response
	return &result, nil
}

type EngagementResponse struct {
	Engagement struct {
		ID                     int    `json:"id,omitempty"`
		PortalID               int    `json:"portalId,omitempty"`
		Active                 bool   `json:"active,omitempty"`
		CreatedAt              int    `json:"createdAt,omitempty"`
		LastUpdated            int    `json:"lastUpdated,omitempty"`
		Type                   string `json:"type,omitempty"`
		UID                    string `json:"uid,omitempty"`
		Timestamp              int    `json:"timestamp,omitempty"`
		Source                 string `json:"source,omitempty"`
		AllAccessibleTeamIds   []any  `json:"allAccessibleTeamIds,omitempty"`
		BodyPreview            string `json:"bodyPreview,omitempty"`
		QueueMembershipIds     []any  `json:"queueMembershipIds,omitempty"`
		BodyPreviewIsTruncated bool   `json:"bodyPreviewIsTruncated,omitempty"`
		BodyPreviewHTML        string `json:"bodyPreviewHtml,omitempty"`
	} `json:"engagement,omitempty"`
	Associations struct {
		TicketIds []int `json:"ticketIds,omitempty"`
	} `json:"associations,omitempty"`
	Attachments []struct {
		ID int `json:"id,omitempty"`
	} `json:"attachments,omitempty"`
	Metadata struct {
		Subject                    string `json:"subject,omitempty"`
		HTML                       string `json:"html,omitempty"`
		Text                       string `json:"text,omitempty"`
		AssociatedContactID        int    `json:"associatedContactId,omitempty"`
		IncomingEmailIsOutOfOffice bool   `json:"incomingEmailIsOutOfOffice,omitempty"`
	} `json:"metadata,omitempty"`
}
