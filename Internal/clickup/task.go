package clickup

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

/*
UpdateTask updates a gives ClickUp task with properties that can be provided in the payload.
It only touches the properties that are given, meaning the status will not be effected if you did not provide it in the payload
*/
func (s clickUpClient) UpdateTask(ctx context.Context, taskID string, payload UpdateTaskPayload) error {

	//generate URL
	endPoint := fmt.Sprintf(ClickupTaskEndpoint, taskID)

	//marshall body
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling commentBody: %w", err)
	}
	// we do not care about the response
	_, err = s.call(ctx, http.MethodPut, endPoint, body)
	if err != nil {
		return fmt.Errorf("error while calling api: %w", err)

	}
	return nil
}

/*
CreateTask calls the ClickUp client with a request to create a task in ClickUp and parses/returns the response.
*/
func (s clickUpClient) CreateTask(ctx context.Context, reqBody CreateTaskRequest) (*CreateTaskResponse, error) {
	//marshall data
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshalling downstream request: %w", err)
	}
	resp, err := s.call(ctx, http.MethodPost, createTaskEndPoint, jsonData)
	if err != nil {
		return nil, fmt.Errorf("error calling api: %w", err)
	}
	//parse response
	var result CreateTaskResponse
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing createTaskResponse: %v", err)
	}
	return &result, nil
}

func (s clickUpClient) GetTask(ctx context.Context, taskID string) (*ClickUpTask, error) {
	endpoint := fmt.Sprintf(ClickupTaskEndpoint, taskID)
	resp, err := s.call(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error calling clickup API: %w", err)
	}
	var result ClickUpTask
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling clickup task response: %w", err)
	}
	return &result, nil
}

// createNewTaskRequest is what we send to the create_task API.
type CreateTaskRequest struct {
	Description  string               `json:"description"` //the ticket FileID from HubSpot
	TicketName   string               `json:"name"`
	CustomFields []ClickupCustomField `json:"custom_fields,omitempty"` //the custom_field customerID in clickup
}

// createTaskResponse is what we expect back from the "create_task" API.
type CreateTaskResponse struct {
	Status ClickUpstatus `json:"status"`
	Name   string        `json:"name"`
	TaskId string        `json:"id"`
}

// status represents the status data structure from clickup in the createTaskResponse struct
type ClickUpstatus struct {
	Id         string `json:"id"`
	Status     string `json:"status"`
	Color      string `json:"color"`
	OrderIndex int    `json:"orderindex"`
	Type       string `json:"type"`
}

//represents what we send to the ClickUp Comment endpoint
type LeaveCommentRequest struct {
	NotifyAll   bool   `json:"notify_all"`
	CommentText string `json:"comment_text"`
}

//represents what we expect from the ClickUp Comment endpoint
type LeaveCommentResponse struct {
	ID     int64  `json:"id"`
	HistID string `json:"hist_id"`
	Date   int64  `json:"date"`
}

//represents what we send to the CLickUp task api when updating
type UpdateTaskPayload struct {
	Status          string `json:"status,omitempty"`
	MarkdownContent string `json:"markdown_content,omitempty"`
}

type ClickUpTask struct {
	ID           string `json:"id"`
	CustomID     any    `json:"custom_id"`
	CustomItemID int    `json:"custom_item_id"`
	Name         string `json:"name"`
	TextContent  string `json:"text_content"`
	Description  string `json:"description"`
	Status       struct {
		ID         string `json:"id"`
		Status     string `json:"status"`
		Color      string `json:"color"`
		Orderindex int    `json:"orderindex"`
		Type       string `json:"type"`
	} `json:"status"`
	Orderindex  string `json:"orderindex"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
	DateClosed  any    `json:"date_closed"`
	DateDone    any    `json:"date_done"`
	Archived    bool   `json:"archived"`
	Creator     struct {
		ID             int    `json:"id"`
		Username       string `json:"username"`
		Color          string `json:"color"`
		Email          string `json:"email"`
		ProfilePicture any    `json:"profilePicture"`
	} `json:"creator"`
	Assignees      []any `json:"assignees"`
	GroupAssignees []any `json:"group_assignees"`
	Watchers       []struct {
		ID             int    `json:"id"`
		Username       string `json:"username"`
		Color          string `json:"color"`
		Initials       string `json:"initials"`
		Email          string `json:"email"`
		ProfilePicture any    `json:"profilePicture"`
	} `json:"watchers"`
	Checklists     []any `json:"checklists"`
	Tags           []any `json:"tags"`
	Parent         any   `json:"parent"`
	TopLevelParent any   `json:"top_level_parent"`
	Priority       any   `json:"priority"`
	DueDate        any   `json:"due_date"`
	StartDate      any   `json:"start_date"`
	Points         any   `json:"points"`
	TimeEstimate   any   `json:"time_estimate"`
	TimeSpent      int   `json:"time_spent"`
	CustomFields   []struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Type       string `json:"type"`
		TypeConfig struct {
		} `json:"type_config,omitempty"`
		DateCreated    string     `json:"date_created"`
		HideFromGuests bool       `json:"hide_from_guests"`
		Required       bool       `json:"required"`
		Value          FlexString `json:"value,omitempty"`
		ValueRichtext  string     `json:"value_richtext,omitempty"`
		TypeConfig0    struct {
			Sorting     string `json:"sorting"`
			NewDropDown bool   `json:"new_drop_down"`
			Options     []struct {
				ID         string `json:"id"`
				Name       string `json:"name"`
				Color      string `json:"color"`
				Orderindex int    `json:"orderindex"`
			} `json:"options"`
		} `json:"type_config,omitempty"`
	} `json:"custom_fields"`
	Dependencies []any  `json:"dependencies"`
	LinkedTasks  []any  `json:"linked_tasks"`
	Locations    []any  `json:"locations"`
	TeamID       string `json:"team_id"`
	URL          string `json:"url"`
	Sharing      struct {
		Public               bool     `json:"public"`
		PublicShareExpiresOn any      `json:"public_share_expires_on"`
		PublicFields         []string `json:"public_fields"`
		Token                any      `json:"token"`
		SeoOptimized         bool     `json:"seo_optimized"`
	} `json:"sharing"`
	PermissionLevel string `json:"permission_level"`
	List            struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Access bool   `json:"access"`
	} `json:"list"`
	Project struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Hidden bool   `json:"hidden"`
		Access bool   `json:"access"`
	} `json:"project"`
	Folder struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Hidden bool   `json:"hidden"`
		Access bool   `json:"access"`
	} `json:"folder"`
	Space struct {
		ID string `json:"id"`
	} `json:"space"`
	Attachments []struct {
		ID               string `json:"id"`
		Date             string `json:"date"`
		Title            string `json:"title"`
		Type             int    `json:"type"`
		Source           int    `json:"source"`
		Version          int    `json:"version"`
		Extension        string `json:"extension"`
		ThumbnailSmall   any    `json:"thumbnail_small"`
		ThumbnailMedium  any    `json:"thumbnail_medium"`
		ThumbnailLarge   any    `json:"thumbnail_large"`
		IsFolder         any    `json:"is_folder"`
		Mimetype         string `json:"mimetype"`
		Hidden           bool   `json:"hidden"`
		ParentID         string `json:"parent_id"`
		Size             int    `json:"size"`
		TotalComments    int    `json:"total_comments"`
		ResolvedComments int    `json:"resolved_comments"`
		User             struct {
			ID             int    `json:"id"`
			Username       string `json:"username"`
			Email          string `json:"email"`
			Initials       string `json:"initials"`
			Color          string `json:"color"`
			ProfilePicture any    `json:"profilePicture"`
		} `json:"user"`
		Deleted     bool   `json:"deleted"`
		Orientation any    `json:"orientation"`
		URL         string `json:"url"`
		EmailData   any    `json:"email_data"`
		WorkspaceID any    `json:"workspace_id"`
		URLWQuery   string `json:"url_w_query"`
		URLWHost    string `json:"url_w_host"`
	} `json:"attachments"`
}

type FlexString string

func (f *FlexString) UnmarshalJSON(data []byte) error {
	// Try string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*f = FlexString(s)
		return nil
	}
	// Try number
	var n json.Number
	if err := json.Unmarshal(data, &n); err == nil {
		*f = FlexString(n.String())
		return nil
	}
	return fmt.Errorf("cannot unmarshal %s into FlexString", data)
}
