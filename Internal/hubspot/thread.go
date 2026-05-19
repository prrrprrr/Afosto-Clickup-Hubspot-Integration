package hubspot

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

/*
GetTicketThreads is a function that returns all threads associated with a HubSpot ticket
*/
func (s hubSpotClient) GetTicketThreads(ctx context.Context, ticketID string) (*HubSpotTicketThreadResponse, error) {
	//prepare endpoint
	endpoint := fmt.Sprintf(HubSpotGetTicketThreadsEndpoint, ticketID)

	//prepare request
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating http request: %w", err)
	}

	//set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.hubSpotApiKey))

	//send request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error executing http request: %w", err)
	}
	defer func() { resp.Body.Close() }()

	//check response code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HubSpot thread-API returned status code: %d, body %s", resp.StatusCode, string((errorBody)))
	}

	//parse response
	var result HubSpotTicketThreadResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("ERROR parsing ticketThreadResponse: %w", err)
	}

	//return parsed response
	return &result, nil
}

/*
PostToTicketThread sends a message to a customer through a HubSpot thread.
*/
func (s hubSpotClient) PostToTicketThread(ctx context.Context, threadID string, request *HubSpotSendEmailRequest) (*HubSpotSendEmailResponse, error) {
	//prepare url
	endPoint := fmt.Sprintf(HubSpotPostToThreadEndpoint, threadID)

	//compose body
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request that send email in thread: %v", err)
	}

	//do request
	resp, err := s.call(ctx, http.MethodPost, endPoint, jsonData)
	if err != nil {
		return nil, fmt.Errorf("error calling thread API: %w", err)
	}

	//parse response
	var result HubSpotSendEmailResponse
	json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing response from engagement API: %w", err)
	}

	//return parsed response
	return &result, nil
}

type HubSpotTicketThreadResponse struct {
	Threads []Thread `json:"results,omitempty"`
}
type Thread struct {
	ID                             string    `json:"id"`
	CreatedAt                      time.Time `json:"createdAt,omitempty"`
	Status                         string    `json:"status,omitempty"`
	OriginalChannelID              string    `json:"originalChannelId,omitempty"`
	OriginalChannelAccountID       string    `json:"originalChannelAccountId,omitempty"`
	LatestMessageTimestamp         time.Time `json:"latestMessageTimestamp,omitempty"`
	LatestMessageReceivedTimestamp time.Time `json:"latestMessageReceivedTimestamp,omitempty"`
	Spam                           bool      `json:"spam,omitempty"`
	InboxID                        string    `json:"inboxId,omitempty"`
	AssociatedContactID            string    `json:"associatedContactId"`
	Archived                       bool      `json:"archived,omitempty"`
}
type HubSpotSendEmailRequest struct {
	Type        string `json:"type,omitempty"`
	Text        string `json:"text,omitempty"`
	Attachments []struct {
		FileID string `json:"fileId,omitempty"`
		Type   string `json:"type,omitempty"`
	} `json:"attachments,omitempty"`
	Recipients       []Recipients `json:"recipients,omitempty"`
	SenderActorID    string       `json:"senderActorId,omitempty"`
	ChannelID        string       `json:"channelId,omitempty"`
	ChannelAccountID string       `json:"channelAccountId,omitempty"`
	RichText         string       `json:"richText,omitempty"`
	Subject          string       `json:"subject,omitempty"`
}
type Attachments struct {
	FileID string `json:"fileId,omitempty"`
	Type   string `json:"type,omitempty"`
}
type DeliveryIdentifier struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}
type Recipients struct {
	ActorID            string             `json:"actorId,omitempty"`
	DeliveryIdentifier DeliveryIdentifier `json:"deliveryIdentifier,omitempty"`
	Name               string             `json:"name,omitempty"`
	RecipientField     string             `json:"recipientField,omitempty"`
}
type HubSpotSendEmailResponse struct {
	ID                    string       `json:"id,omitempty"`
	ConversationsThreadID string       `json:"conversationsThreadId,omitempty"`
	CreatedAt             time.Time    `json:"createdAt,omitempty"`
	UpdatedAt             time.Time    `json:"updatedAt,omitempty"`
	CreatedBy             string       `json:"createdBy,omitempty"`
	Client                Client       `json:"client,omitempty"`
	Senders               []Senders    `json:"senders,omitempty"`
	Recipients            []Recipients `json:"recipients,omitempty"`
	Archived              bool         `json:"archived,omitempty"`
	Text                  string       `json:"text,omitempty"`
	RichText              string       `json:"richText,omitempty"`
	Attachments           []any        `json:"attachments,omitempty"`
	Subject               string       `json:"subject,omitempty"`
	TruncationStatus      string       `json:"truncationStatus,omitempty"`
	Status                Status       `json:"status,omitempty"`
	Direction             string       `json:"direction,omitempty"`
	ChannelID             string       `json:"channelId,omitempty"`
	ChannelAccountID      string       `json:"channelAccountId,omitempty"`
	Type                  string       `json:"type,omitempty"`
}
type Client struct {
	ClientType       string `json:"clientType,omitempty"`
	IntegrationAppID int    `json:"integrationAppId,omitempty"`
}
type Senders struct {
	ActorID            string             `json:"actorId,omitempty"`
	Name               string             `json:"name,omitempty"`
	SenderField        string             `json:"senderField,omitempty"`
	DeliveryIdentifier DeliveryIdentifier `json:"deliveryIdentifier,omitempty"`
}
type Status struct {
	StatusType string `json:"statusType,omitempty"`
}
