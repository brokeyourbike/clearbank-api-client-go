package clearbank

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	defaultBaseURL = "https://institution-api-sim.clearbank.co.uk"
)

// requestIdCtx is the context key for the request ID.
type requestIdCtx struct{}

// Client is a ClearBank API client.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
	signer     Signer
}

// ClientOption is a function that configures a Client.
type ClientOption func(*Client)

func NewClient(token string, signer Signer, options ...ClientOption) *Client {
	c := &Client{
		baseURL:    defaultBaseURL,
		token:      token,
		httpClient: http.DefaultClient,
		signer:     signer,
	}

	for _, option := range options {
		option(c)
	}

	return c
}

// WithHTTPClient sets the HTTP client for the ClearBank API client.
func WithHTTPClient(c *http.Client) ClientOption {
	return func(client *Client) {
		client.httpClient = c
	}
}

// WithBaseURL sets the base URL for the ClearBank API client.
func WithBaseURL(baseURL string) ClientOption {
	return func(client *Client) {
		client.baseURL = strings.TrimSuffix(baseURL, "/")
	}
}

func (c *Client) NewRequest(ctx context.Context, method, url string, body interface{}) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if req.Method == http.MethodPost || req.Method == http.MethodPatch {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}

		signature, err := c.signer.Sign(ctx, b)
		if err != nil {
			return nil, fmt.Errorf("failed to sign payload: %w", err)
		}

		req.Body = io.NopCloser(bytes.NewReader(b))
		req.Header.Set("DigitalSignature", string(signature))
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-Id", RequestIdFromContext(ctx))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return resp, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	return resp, nil
}

// RequestIdContext adds a request ID to the given context.
func RequestIdContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIdCtx{}, id)
}

// RequestIdFromContext returns the request ID from the given context.
func RequestIdFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if id, ok := ctx.Value(requestIdCtx{}).(string); ok {
		return id
	}

	return ""
}
