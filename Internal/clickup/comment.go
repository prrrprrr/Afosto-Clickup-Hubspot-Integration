package clickup

/*
comment.go contains the client's methods that relate to comments.
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

/*
LeaveTaskComment creates a comment with the HubSpot ticket link in the newly created clickup task.
We use a comment because ClickUp creates a relationship with the HubSpot Ticket that way and stores the URL.
*/
func (s clickUpClient) LeaveTaskComment(ctx context.Context, taskId string, request LeaveCommentRequest) (*LeaveCommentResponse, error) {
	//compose endpoint
	endpoint := fmt.Sprintf(clickupCommentEndpoint, taskId)

	//marshall body
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling commentBody: %w", err)
	}

	//do request
	resp, err := s.call(ctx, http.MethodPost, endpoint, jsonData)

	//parse response
	var result LeaveCommentResponse
	if err = json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error parsing leaveCommentResponse: %w", err)
	}

	//return response
	return &result, nil
}

/*
LeaveThreadedComment leaves a threaded comment under an already existing comment. The comment to be threaded is in the request and the comment to thread under is gives as commentID.
*/
func (s clickUpClient) LeaveThreadedComment(ctx context.Context, commentID string, request LeaveCommentRequest) (*LeaveCommentResponse, error) {
	//compose endpoint
	endpoint := fmt.Sprintf(ThreadedCommentEndpoint, commentID)

	//marshall body
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling commentBody: %w", err)
	}

	//do request
	resp, err := s.call(ctx, http.MethodPost, endpoint, jsonData)

	//parse response
	var result LeaveCommentResponse
	if err = json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error parsing leaveCommentResponse: %w", err)
	}

	//return response
	return &result, nil
}

/*
GetComments retrieves all the comments of a clickup task and returns them
*/
func (s clickUpClient) GetComments(ctx context.Context, taskId string) (*getCommentResponse, error) {
	//compose endpoint
	endpoint := fmt.Sprintf(clickupCommentEndpoint, taskId)

	//do request
	resp, err := s.call(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error sending getComments request %w", err)
	}

	//parse resp
	var result getCommentResponse
	if err = json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error parsing getCommentResponse: %w", err)
	}

	//return resp
	return &result, nil
}

/*
parseClickUpComment uses the blocksToHTML helper function to parse clickup responses into usable richtext for outbound emails.
*/
func (s clickUpClient) parseClickUpComment(ctx context.Context, response *Comment) (string, error) {
	//TODO error handling in this function and BlocksToHTML
	newHTML := BlocksToHTML(response.Comment)
	return newHTML, nil
}

/*
GetLastComment exclusively returns the most recent clickup comment
it returns it already parsed into HTML so the caller will not have to parse it.
*/
func (s clickUpClient) GetLastComment(ctx context.Context, taskID string) (*ParsedComment, error) {
	//getting comments
	allComments, err := s.GetComments(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("error getting comments: %w", err)
	}

	//getting the parsed version of the comment
	commentText, err := s.parseClickUpComment(ctx, &allComments.Comments[0])

	//prepare response
	parsedComment := ParsedComment{
		Comment: commentText,
		ID:      allComments.Comments[0].ID,
		User:    allComments.Comments[0].User.Email,
	}

	//return response
	return &parsedComment, nil
}

type ParsedComment struct {
	Comment string
	ID      string
	User    string
}

type getCommentResponse struct {
	Comments []Comment `json:"comments,omitempty"`
}

type Comment struct {
	ID          string         `json:"id,omitempty"`
	Comment     []CommentBlock `json:"comment,omitempty"`
	CommentText string         `json:"comment_text,omitempty"`
	User        struct {
		ID             int    `json:"id,omitempty"`
		Username       string `json:"username,omitempty"`
		Email          string `json:"email,omitempty"`
		Color          string `json:"color,omitempty"`
		Initials       string `json:"initials,omitempty"`
		ProfilePicture string `json:"profilePicture,omitempty"`
	} `json:"user,omitempty"`
	Assignee      any    `json:"assignee,omitempty"`
	GroupAssignee any    `json:"group_assignee,omitempty"`
	Reactions     []any  `json:"reactions,omitempty"`
	Date          string `json:"date,omitempty"`
	ReplyCount    int    `json:"reply_count,omitempty"`
}

type CommentBlock struct {
	Text        string          `json:"text"`
	Type        string          `json:"type"` // "emoticon", "tag", or empty
	Attributes  BlockAttrs      `json:"attributes"`
	Emoticon    *EmoticonData   `json:"emoticon"`
	User        *UserRef        `json:"user"`
	Image       *ImageData      `json:"image"`
	LinkMention *LinkMention    `json:"link_mention"`
	Attachment  *AttachmentData `json:"attachment"`
	Bookmark    *BookmarkData   `json:"bookmark"`
}
type BookmarkData struct {
	Service string `json:"service"`
	ID      string `json:"id"`
	Url     string `json:"url"`
}

type BlockAttrs struct {
	Bold       bool       `json:"bold"`
	Italic     bool       `json:"italic"`
	Strike     bool       `json:"strike"`
	Underline  bool       `json:"underline"`
	Code       bool       `json:"code"`
	Link       string     `json:"link"`
	Color      string     `json:"color"`
	Background string     `json:"background"`
	Indent     int        `json:"indent"`
	List       *ListAttr  `json:"list"`
	CodeBlock  *CodeBlock `json:"code-block"`
	BlockID    string     `json:"block-id"` // Added to handle ClickUp's block-id
	Width      string     `json:"width"`
}

type ImageData struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

type AttachmentData struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	Size int64  `json:"size"` // bytes; may be 0 if not provided
}

type ListAttr struct {
	List string `json:"list"`
}

type CodeBlock struct {
	Lang string
}

type EmoticonData struct {
	Code string `json:"code"` // hex code point e.g. "1f60a"
}

type LinkMention struct {
	URL string `json:"url"`
}

type UserRef struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
