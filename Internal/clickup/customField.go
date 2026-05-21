package clickup

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

/*
GetCustomerID has the job of searching the customer custom field in ClickUp to see if there is a match with payload.companyname.
If we find a match we fill in the custom field with the found company, if not we use a fallback company so it is clear we have not found a match.
*/
func (s clickUpClient) GetCustomerID(ctx context.Context, customer string) (string, error) {
	//call api to collect current custom_id information
	resp, err := s.GetCustomFields(ctx)
	if err != nil {
		return "", fmt.Errorf("error fetching custom fields from list: %w", err)
	}
	//loop over all custom fields until 'customer' has been found
	for _, v := range resp.Fields {
		if v.Name == "Customer" {
			//loop over all possible entries and look for a match with payload.companyname
			for _, b := range v.TypeConfig.Options {
				//if we find it we return the result
				if b.Name == customer {
					return b.ID, nil
				}
			}
		}
	}
	//return an error if we could not find the correct customer
	return "", fmt.Errorf("error finding custom field FileID for company name")
}

/*
GetCustomFields makes a GET request to the custom field endpoint to
retrieve the custom fields related to the given ClickUp list
*/
func (s clickUpClient) GetCustomFields(ctx context.Context) (*CustomFieldResponse, error) {

	resp, err := s.call(ctx, http.MethodGet, getCustomFieldEndPoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting custom fields:%s", err)

	}

	//parse response
	var response CustomFieldResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return nil, fmt.Errorf("error parsing CustomFieldResponse")
	}

	//return response
	return &response, nil
}

// customField represents a custom field in Clickup in the createTaskResponse
type ClickupCustomField struct {
	CustomFieldID    string `json:"id"`
	CustomFieldValue string `json:"value"`
}

// customFieldResponse is what we expect back from the custom_field api
type CustomFieldResponse struct {
	Fields []struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Type       string `json:"type"`
		TypeConfig struct {
			Sorting     string `json:"sorting"`
			NewDropDown bool   `json:"new_drop_down"`
			Options     []struct {
				ID         string `json:"id"`
				Name       string `json:"name"`
				Color      string `json:"color"`
				OrderIndex int    `json:"order_index"`
			} `json:"options"`
		} `json:"type_config"`
		DateCreated    string `json:"date_created"`
		HideFromGuests bool   `json:"hide_from_guests"`
		Required       bool   `json:"required"`
	} `json:"fields"`
}
