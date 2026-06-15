package clickup

/*
This file contains client functions regarding sprints
*/
import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s clickUpClient) MoveTaskToCurrentSprint(ctx context.Context, taskID string) (*MoveTaskResponse, error) {
	//get current sprint
	sprintID, err := s.getCurrentSprintID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error moving task to current sprint: %w", err)
	}

	//compose endpoint
	endpoint := fmt.Sprintf(moveTaskEndpoint, taskID, sprintID)

	//call enpoint
	resp, err := s.call(ctx, http.MethodPut, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error moving the task to the current sprint: %w", err)
	}
	//parse response
	var moveTaskResp MoveTaskResponse
	err = json.Unmarshal(resp, &moveTaskResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling the response for moving a task to a different list in clickup: %w", err)
	}
	//return response
	return &moveTaskResp, nil
}

//Helper function for moving a task to the current sprint. Since the sprintID changes every two weeks we need to pull it everytime.
func (s clickUpClient) getCurrentSprintID(ctx context.Context) (string, error) {
	//compose endpoint
	endpoint := fmt.Sprintf(getListsEndpoint, sprintFolderID)
	//call request
	resp, err := s.call(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("error: could not get current sprint id from clickup API: %w", err)
	}

	//parse response
	var getListsReponse GetListsResponse
	err = json.Unmarshal(resp, &getListsReponse)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling the response for getting lists in a folder: %w", err)
	}

	//return response.
	return getListsReponse.Lists[0].ID, nil
}

type MoveTaskResponse struct {
	Data struct {
		TaskID    string `json:"task_id"`
		NewListId string `json:"new_list_id"`
	}
}

type GetListsResponse struct {
	Lists []Lists `json:"lists,omitempty"`
}
type Folder struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Hidden   bool   `json:"hidden,omitempty"`
	Archived bool   `json:"archived,omitempty"`
	Access   bool   `json:"access,omitempty"`
}
type Space struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Archived bool   `json:"archived,omitempty"`
	Access   bool   `json:"access,omitempty"`
}
type Lists struct {
	ID               string `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	Orderindex       int    `json:"orderindex,omitempty"`
	Content          string `json:"content,omitempty"`
	Status           any    `json:"status,omitempty"`
	Priority         any    `json:"priority,omitempty"`
	Assignee         any    `json:"assignee,omitempty"`
	TaskCount        int    `json:"task_count,omitempty"`
	DueDate          string `json:"due_date,omitempty"`
	StartDate        string `json:"start_date,omitempty"`
	Folder           Folder `json:"folder,omitempty"`
	Space            Space  `json:"space,omitempty"`
	Archived         bool   `json:"archived,omitempty"`
	OverrideStatuses bool   `json:"override_statuses,omitempty"`
	PermissionLevel  string `json:"permission_level,omitempty"`
}
