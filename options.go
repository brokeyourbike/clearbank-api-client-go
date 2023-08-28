package clearbank

import (
	"strings"
)

// WithHTTPClient sets the HTTP client for the ClearBank API client.
func WithHTTPClient(c HttpClient) ClientOption {
	return func(target *client) {
		target.httpClient = c
	}
}

// WithBaseURL sets the base URL for the ClearBank API client.
func WithBaseURL(baseURL string) ClientOption {
	return func(target *client) {
		target.baseURL = strings.TrimSuffix(baseURL, "/")
	}
}
