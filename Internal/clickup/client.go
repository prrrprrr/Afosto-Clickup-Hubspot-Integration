package clickup

/*
client.go contains the interface, settings, factoryMethod and call function for the clickup client.
*/
import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

type ClickUpClient interface {
	CreateTask(ctx context.Context, payload CreateTaskRequest) (response *CreateTaskResponse, err error)
	UpdateTask(ctx context.Context, taskID string, payload UpdateTaskPayload) error
	GetTask(ctx context.Context, taskID string) (task *ClickUpTask, err error)
	GetCustomerID(ctx context.Context, customer string) (response string, err error)
	LeaveTaskComment(ctx context.Context, taskID string, request LeaveCommentRequest) (response *LeaveCommentResponse, err error)
	LeaveThreadedComment(ctx context.Context, taskID string, request LeaveCommentRequest) (*LeaveCommentResponse, error)
	GetComments(ctx context.Context, taskID string) (response *getCommentResponse, err error)
	GetLastComment(ctx context.Context, taskID string) (*ParsedComment, error)
	GetCustomFields(ctx context.Context) (response *CustomFieldResponse, err error)
	UploadFileFromUrl(ctx context.Context, file HubspotFile, taskID string) (response *UploadImageResponse, err error)
	MoveTaskToCurrentSprint(ctx context.Context, taskID string) (resp *MoveTaskResponse, err error)
	getCurrentSprintID(ctx context.Context) (sprintID string, err error)
}

type Settings struct {
	Client        *http.Client
	ClickupApiKey string
	HubSpotAPiKey string
}

type clickUpClient struct {
	client        *http.Client
	clickupApiKey string
	HubSpotAPiKey string
}

func NewClient(settings Settings) ClickUpClient {

	return clickUpClient{
		client:        settings.Client,
		clickupApiKey: settings.ClickupApiKey,
		HubSpotAPiKey: settings.HubSpotAPiKey,
	}
}

func (s clickUpClient) call(ctx context.Context, method string, endpoint string, payload []byte) ([]byte, error) {
	req, err := http.NewRequest(method, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating httpRequest: %w", err)
	}

	//set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.clickupApiKey)
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
		return nil, fmt.Errorf("ClickUp API returned status code: %d, body %s", resp.StatusCode, string((body)))
	}
	//result response
	return body, nil
}
