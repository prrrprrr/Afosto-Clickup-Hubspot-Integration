package hubspot

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

type HubSpotClient interface {
	UpdateTicket(ctx context.Context, ticketID string, request UpdateTicketRequest) (response *UpdateTicketResponse, err error)
	GetTicketThreads(ctx context.Context, ticketID string) (response *HubSpotTicketThreadResponse, err error)
	PostToTicketThread(ctx context.Context, threadID string, request *HubSpotSendEmailRequest) (response *HubSpotSendEmailResponse, err error)
	UploadFileFromUrl(ctx context.Context, ticketID string, url string) (response *FileResponse, err error)
	GetFile(ctx context.Context, fileID string) (response *FileResponse, err error)
	GetEngagement(ctx context.Context, engagementID string) (response *EngagementResponse, err error)
	GetSignedUrl(ctx context.Context, fileID string) (*SignedUrl, error)
	StatusToPipeLine(status string) (string, error)
}
type Settings struct {
	Client *http.Client
	APIKey string
}
type hubSpotClient struct {
	client        *http.Client
	hubSpotApiKey string
}

func NewClient(settings Settings) HubSpotClient {
	return hubSpotClient{
		client:        settings.Client,
		hubSpotApiKey: settings.APIKey,
	}
}

func (s hubSpotClient) call(ctx context.Context, method string, endpoint string, payload []byte) ([]byte, error) {
	req, err := http.NewRequest(method, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating httpRequest: %w", err)
	}

	//set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.hubSpotApiKey))
	req.Header.Set("accept", "application/json")

	//call api
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	// Read the full body now
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	//check response code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Hubspot API returned status code: %d, with body: %s, from endpoint: %s", resp.StatusCode, string(body), endpoint)
	}
	//result response
	return body, nil
}
