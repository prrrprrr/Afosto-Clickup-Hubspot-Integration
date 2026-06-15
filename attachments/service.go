package attachments

/*
attachmentsService.go contains the interface for the attachmentService and all functions related to it
*/

import (
	"Afosto-Clickup-Hubspot-Integration/internal/clickup"
	"Afosto-Clickup-Hubspot-Integration/internal/hubspot"
	"context"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type AttachmentService interface {
	findInlineImages(ctx context.Context, htmlStr string) (sources []InlineImage, err error)
	dfsSources(node *html.Node, sources *[]InlineImage)
	replaceImageSources(ctx context.Context, htmlStr string, images []clickup.ClickupFile) (newString string, err error)
	UploadImagesToClickUpTask(ctx context.Context, taskID string, images *[]clickup.HubspotFile) (resp []clickup.ClickupFile, err error)
	HandleInlineAttachments(ctx context.Context, taskID string, engagementID string) (err error)
	HandleAttachments(ctx context.Context, taskID string, engagement hubspot.EngagementResponse) (err error)
	ProcessAttachments(ctx context.Context, taskID string, engagement hubspot.EngagementResponse) (err error)
}
type attachmentService struct {
	ClickupClient clickup.ClickUpClient
	HubspotClient hubspot.HubSpotClient
}

func NewService(
	clickupClient clickup.ClickUpClient,
	hubspotClient hubspot.HubSpotClient,
) AttachmentService {
	return &attachmentService{
		ClickupClient: clickupClient,
		HubspotClient: hubspotClient,
	}

}

/*
findInlineImages checks the incoming mail for any images in the htmlStr. It has a helper function named dfsSources
If there are images, the sources of these images are linked to the OriginalSource's and returned as an InlineImage slice.
*/
func (s *attachmentService) findInlineImages(ctx context.Context, htmlStr string) ([]InlineImage, error) {
	if len(htmlStr) <= 0 {
		return nil, fmt.Errorf("error the htmlString was empty")
	}
	doc, _ := html.Parse(strings.NewReader(htmlStr))
	var sources []InlineImage
	s.dfsSources(doc, &sources)
	return sources, nil
}

/*
dfsSources collects the id and cid of hubspot images and returns them as an InlineImage slice
dfs stands for Depth First Search and is the algorithm used.
*/
func (s *attachmentService) dfsSources(node *html.Node, sources *[]InlineImage) {
	//we take the
	if node.Type == html.ElementNode && node.Data == "img" {
		originalSource := ""
		id := ""
		url := ""
		for _, attr := range node.Attr {
			if attr.Key == "data-file-id" {
				id = attr.Val
			}
			if attr.Key == "data-original-src" {
				originalSource = attr.Val
			}
			if attr.Key == "src" {
				url = attr.Val
			}
		}
		*sources = append(*sources, InlineImage{OriginalSource: originalSource, FileID: id, URL: url})
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		s.dfsSources(c, sources)
	}
}

/*
UploadImagesToClickUpTask uploads the images to the clickup Task so they're available in the task and usable in clickup
It returns a pointer to a []ImageSource that contains the new CLickUpURL and the already existing OriginalSource
*/
func (s *attachmentService) UploadImagesToClickUpTask(ctx context.Context, taskID string, images *[]clickup.HubspotFile) ([]clickup.ClickupFile, error) {
	var results []clickup.ClickupFile
	for _, img := range *images {
		resp, err := s.ClickupClient.UploadFileFromUrl(ctx, img, taskID)
		if err != nil {
			return nil, fmt.Errorf("could not upload image: %w", err)
		}
		results = append(results, clickup.ClickupFile{URL: resp.URL, CID: img.CID, FileName: resp.Name, Extension: resp.Extension})
	}
	return results, nil
}

/*
replaceImageSources replaces the original sources of the description of a ticket with the new clickup sources so the images are visible in the description
This is not used anywhere as the requirement fell through, however it is still here
*/
func (s *attachmentService) replaceImageSources(ctx context.Context, description string, sources []clickup.ClickupFile) (string, error) {
	//go through text
	//replace hubspot sources with new clickup links

	for _, source := range sources {
		cid := fmt.Sprintf("[%s]", source.CID)
		if source.CID == "" {
			continue // skip images with no OriginalSource — replacing "" would corrupt the entire string
		}
		url := fmt.Sprintf("![image](%s)", source.URL)
		description = strings.Replace(description, cid, url, 1)
	}
	//parse string to find OriginalSource's
	return description, nil
}

/*
	HandleInlineAttachments is a  function that handles the inline attachments from incoming emails
	Asks for a htmlString and the ClickUpTaskId so it can extract the images and upload them to the task
	HandleAttachments is for handling all classic attached files. This function does not cover that.
*/
func (s *attachmentService) HandleInlineAttachments(ctx context.Context, taskID string, htmlString string) error {
	var images []clickup.HubspotFile
	//find all images in the email body
	sources, err := s.findInlineImages(ctx, htmlString)
	if err != nil {
		return fmt.Errorf("error finding images: %w", err)
	}

	//if we found nothing we can exit early
	if len(sources) == 0 {
		return nil
	}

	//set up the images for processing
	var file clickup.HubspotFile
	for _, source := range sources {
		//get data from hubspot api
		if strings.HasPrefix(source.OriginalSource, "cid") {
			signedURl, err := s.HubspotClient.GetSignedUrl(ctx, source.FileID)
			if err != nil {
				return fmt.Errorf("could not download image: %w", err)
			}
			//save it in a clickup native struct
			file = clickup.HubspotFile{
				FileName:  signedURl.Name,
				URL:       signedURl.URL,
				Extension: signedURl.Extension,
			}
		}
		if strings.HasPrefix(source.OriginalSource, "http") {

			//save it in a clickup native struct
			parsed, err := url.Parse(source.URL)
			if err != nil {
				return fmt.Errorf("could not parse image url: %w", err)
			}
			file = clickup.HubspotFile{
				FileName:  path.Base(parsed.Path),
				URL:       parsed.Query().Get("url"),
				Extension: path.Ext(parsed.Path),
			}
		}

		images = append(images, file)
	}

	//if we found images we upload them to the task
	_, err = s.UploadImagesToClickUpTask(ctx, taskID, &images)
	if err != nil {
		return fmt.Errorf("error uploading images to CLickUp: %w", err)
	}

	return nil
}

/*
	HandleAttachments handles the attached files from an engagement by simply uploading them
	HandleInlineAttachments is for the inline attachments, which this function does not cover
*/
func (s attachmentService) HandleAttachments(ctx context.Context, taskID string, engagement hubspot.EngagementResponse) error {

	if len(engagement.Attachments) < 1 {
		return nil
	}

	for _, attachment := range engagement.Attachments {
		signedURl, err := s.HubspotClient.GetSignedUrl(ctx, strconv.Itoa(attachment.ID))
		//resp, err := s.HubspotClient.GetFile(ctx, signedURl.URL)
		if err != nil {
			return fmt.Errorf("error getting file: %w", err)
		}

		//transfer to clickup struct
		file := clickup.HubspotFile{
			FileName:  signedURl.Name,
			URL:       signedURl.URL,
			Extension: signedURl.Extension,
		}
		_, err = s.ClickupClient.UploadFileFromUrl(ctx, file, taskID)
		if err != nil {
			return fmt.Errorf("error uploading file: %w", err)
		}
	}
	return nil
}

/*
ProcessAttachments uses the other helper functions and handles all attachments of an incoming email. It uploads attached files and also handles the inline attachments
*/
func (s attachmentService) ProcessAttachments(ctx context.Context, taskID string, engagement hubspot.EngagementResponse) error {
	//if there are any attached files
	if len(engagement.Attachments) > 0 {
		err := s.HandleAttachments(ctx, taskID, engagement)
		if err != nil {
			return fmt.Errorf("error handling attachments: %w", err)
		}
	}
	//handle inline images
	err := s.HandleInlineAttachments(ctx, taskID, engagement.Metadata.HTML)
	if err != nil {
		return fmt.Errorf("error handling inline attachments: %w", err)
	}

	return nil

}

type ImageSource struct {
	URL            string
	OriginalSource string
	ID             string
}

type InlineImage struct {
	OriginalSource string
	URL            string
	FileID         string
}
