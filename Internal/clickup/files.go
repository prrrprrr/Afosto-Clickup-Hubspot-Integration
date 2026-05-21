package clickup

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

func (s clickUpClient) UploadFileFromUrl(ctx context.Context, file HubspotFile, taskID string) (*UploadImageResponse, error) {

	imageURL := file.URL
	log.Printf("file.URL is: %s", file.URL)
	// Download the image
	imgResp, err := http.Get(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer imgResp.Body.Close()

	// Build multipart body
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	//compose the filename and extension

	filename := fmt.Sprintf("%s.%s", file.FileName, file.Extension)
	if file.FileName == "" || file.Extension == "" {
		filename = "image.jpeg"
	}
	part, err := writer.CreateFormFile("attachment", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	// Stream image bytes directly into the multipart body
	if _, err = io.Copy(part, imgResp.Body); err != nil {
		return nil, fmt.Errorf("failed to copy image data: %w", err)
	}
	writer.Close()
	// Upload to ClickUp
	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf(attachmentEndpoint, taskID),
		&buf,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", s.clickupApiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	uploadResp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to upload: %w", err)
	}
	defer uploadResp.Body.Close()

	if uploadResp.StatusCode < 200 || uploadResp.StatusCode >= 300 {
		errorBody, _ := io.ReadAll(uploadResp.Body)
		return nil, fmt.Errorf("getting a file from hubspot returned status code: %d, body %s", uploadResp.StatusCode, string((errorBody)))
	}
	uploadBody, err := io.ReadAll(uploadResp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	var response UploadImageResponse
	err = json.Unmarshal(uploadBody, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return &response, nil
}

type UploadImageResponse struct {
	ID              string `json:"id,omitempty"`
	Version         string `json:"version,omitempty"`
	Date            int64  `json:"date,omitempty"`
	Name            string `json:"name,omitempty"`
	Title           string `json:"title,omitempty"`
	Extension       string `json:"extension,omitempty"`
	Source          int    `json:"source,omitempty"`
	ThumbnailSmall  string `json:"thumbnail_small,omitempty"`
	ThumbnailMedium string `json:"thumbnail_medium,omitempty"`
	ThumbnailLarge  string `json:"thumbnail_large,omitempty"`
	Width           int    `json:"width,omitempty"`
	Height          int    `json:"height,omitempty"`
	URL             string `json:"url,omitempty"`
	URLWQuery       string `json:"url_w_query,omitempty"`
	URLWHost        string `json:"url_w_host,omitempty"`
}
type HubspotFile struct {
	FileName  string
	URL       string
	Extension string
	CID       string
}
type ClickupFile struct {
	FileName  string
	URL       string
	Extension string
	CID       string
}
type SignedUrl struct {
	URL       string    `json:"url"`
	ExpiresAt time.Time `json:"expiresAt"`
	Name      string    `json:"name"`
	Extension string    `json:"extension"`
	Type      string    `json:"type"`
	Size      int       `json:"size"`
	Height    int       `json:"height"`
	Width     int       `json:"width"`
}
