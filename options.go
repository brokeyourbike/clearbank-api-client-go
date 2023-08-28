package clearbank

import (
	"strings"

	"github.com/sirupsen/logrus"
)

// WithHTTPClient sets the HTTP client for the ClearBank API client.
func WithHTTPClient(c HttpClient) ClientOption {
	return func(target *client) {
		target.httpClient = c
	}
}

// WithLogger sets the *logrus.Logger for the ClearBank API client.
func WithLogger(l *logrus.Logger) ClientOption {
	return func(target *client) {
		target.logger = l
	}
}

// WithBaseURL sets the base URL for the ClearBank API client.
func WithBaseURL(baseURL string) ClientOption {
	return func(target *client) {
		target.baseURL = strings.TrimSuffix(baseURL, "/")
	}
}
