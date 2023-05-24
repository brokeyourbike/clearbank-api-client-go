package clearbank

import "net/http"

// ErrorResponse represents an error response from the ClearBank API.
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response

	Type     string              `json:"type"`
	Title    string              `json:"title"`
	Status   int                 `json:"status"`
	Details  string              `json:"detail"`
	Instance string              `json:"instance"`
	Errors   map[string][]string `json:"errors"`
}
