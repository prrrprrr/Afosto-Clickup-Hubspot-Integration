package hubspot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"time"
)

/*
GetSignedUrl requests a signed url from the files api.
*/
func (s hubSpotClient) GetSignedUrl(ctx context.Context, fileID string) (*SignedUrl, error) {
	if fileID == "" {
		return nil, fmt.Errorf("couldn't download image because there is no id given")
	}
	endpoint := fmt.Sprintf(HubSpotSignedUrlEndpoint, fileID)
	var signedUrl SignedUrl
	urlResp, err := s.call(ctx, http.MethodGet, endpoint, nil)
	log.Printf("for id: %s, getURL raw response: %s", fileID, string(urlResp))
	if err != nil {
		log.Printf("didn't get results for: %s", fileID)
		return nil, fmt.Errorf("failed to get signed url: %w", err)

	}
	err = json.Unmarshal(urlResp, &signedUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get unmarshal signed url: %w", err)
	}
	if signedUrl.URL == "" {
		log.Printf("signed url is empty")
		return nil, fmt.Errorf("signed url is empty: %w", err)
	}
	return &signedUrl, nil

}

/*
GetFile retrieves a file by its FileID from the HubSpot Files API
*/
func (s hubSpotClient) GetFile(ctx context.Context, fileID string) (*FileResponse, error) {
	//prepare endpoint
	endPoint := fmt.Sprintf(GetFileEndpoint, fileID)

	//compose request
	resp, err := s.call(ctx, http.MethodGet, endPoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error sending request to HubSpot getFile API: %w", err)

	}

	//parse response
	var result FileResponse
	log.Printf("GetFile raw response: %s", string(resp))
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing response to Hubspot struct: %w", err)
	}
	//return parsed response
	return &result, nil
}

/*
UploadFileFromUrl uploads a file to HubSpot
does not use the call method because it is the only function to use Form data
*/
func (s hubSpotClient) UploadFileFromUrl(ctx context.Context, ticketID string, fileURL string) (*FileResponse, error) {

	// Download the image
	imgResp, err := http.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer imgResp.Body.Close()

	// create buffer and writer
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Use the last part of the URL as fileName, fallback to "image.png"
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	fileName := filepath.Base(parsedURL.Path)
	if fileName == "" || fileName == "." {
		fileName = "image.png"
	}
	if filepath.Ext(fileName) == "" {
		fileName += ".png"
	}
	//add file field
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	_, params, err := mime.ParseMediaType(imgResp.Header.Get("content-type"))
	charset := params["charset"]

	// add the other fields
	writer.WriteField("charsetHunch", charset)
	writer.WriteField("fileName", fileName)
	writer.WriteField("folderId", folderID)
	writer.WriteField("options", `{"access": "PUBLIC_NOT_INDEXABLE"}`)

	// Stream image bytes directly into the multipart body
	if _, err = io.Copy(part, imgResp.Body); err != nil {
		return nil, fmt.Errorf("failed to copy image data: %w", err)
	}
	writer.Close()

	//create request
	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf(PostFileEndpoint, ticketID),
		&buf,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	//set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.hubSpotApiKey))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	//do request
	uploadResp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to upload: %w", err)
	}
	//check response code
	if uploadResp.StatusCode < 200 || uploadResp.StatusCode >= 300 {
		errorBody, _ := io.ReadAll(uploadResp.Body)
		return nil, fmt.Errorf("getting a file from hubspot returned status code: %d, body %s", uploadResp.StatusCode, string((errorBody)))
	}
	defer uploadResp.Body.Close()

	uploadBody, err := io.ReadAll(uploadResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	//check response
	var response FileResponse
	err = json.Unmarshal(uploadBody, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	return &response, nil
}

type FileResponse struct {
	Name              string `json:"name,omitempty"`
	Size              int    `json:"size,omitempty"`
	Height            int    `json:"height,omitempty"`
	Width             int    `json:"width,omitempty"`
	Encoding          string `json:"encoding,omitempty"`
	Type              string `json:"type,omitempty"`
	Extension         string `json:"extension,omitempty"`
	DefaultHostingURL string `json:"defaultHostingUrl,omitempty"`
	URL               string `json:"url,omitempty"`
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

type ProxyUrl struct {
	URL string
}
