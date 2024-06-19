package clearbank

import "fmt"

type UnexpectedResponse struct {
	Status int
	Body   string
}

func (r UnexpectedResponse) Error() string {
	return fmt.Sprintf("Unexpected response from API. Status: %d Body: %s", r.Status, r.Body)
}

type ErrResponse struct {
	Type     string              `json:"type"`
	Title    string              `json:"title"`
	Status   int                 `json:"status" validate:"required"`
	Details  string              `json:"detail"`
	Instance string              `json:"instance"`
	Errors   map[string][]string `json:"errors"`
}

func (e ErrResponse) Error() string {
	if len(e.Errors) == 0 {
		return fmt.Sprintf("Error during API call. Status: %d Type: %s Title: %s", e.Status, e.Type, e.Title)
	}
	return fmt.Sprintf("Error during API call. Status: %d Type: %s Title: %s Errors: %s", e.Status, e.Type, e.Title, e.Errors)
}
