package clearbank

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/brokeyourbike/clearbank-api-client-go/signature"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

const defaultBaseURL = "https://institution-api-sim.clearbank.co.uk"

// requestIdCtx is the context key for the request ID.
type requestIdCtx struct{}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client interface {
	TestClient
	RateClient
	StatementClient
}

var _ Client = (*client)(nil)

type client struct {
	httpClient HttpClient
	signer     signature.Signer
	logger     *logrus.Logger
	validate   *validator.Validate
	baseURL    string
	token      string
}

// ClientOption is a function that configures a Client.
type ClientOption func(*client)

func NewClient(token string, signer signature.Signer, options ...ClientOption) *client {
	c := &client{
		httpClient: http.DefaultClient,
		signer:     signer,
		validate:   validator.New(validator.WithRequiredStructEnabled()),
		baseURL:    defaultBaseURL,
		token:      token,
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func (c *client) newRequest(ctx context.Context, method, url string, body interface{}) (*request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var b []byte

	if req.Method == http.MethodPost || req.Method == http.MethodPatch {
		b, err = json.Marshal(body)
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

	requestID := RequestIdFromContext(ctx)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-Id", requestID)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	if c.logger != nil {
		c.logger.WithContext(ctx).WithFields(logrus.Fields{
			"http.request.body.content":       string(b),
			"http.request.headers.request_id": requestID,
		}).Debug("clearbank.client -> request")
	}

	return NewRequest(req), nil
}

func (c *client) do(ctx context.Context, req *request) error {
	resp, err := c.httpClient.Do(req.req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(b))

	if c.logger != nil {
		c.logger.WithContext(ctx).WithFields(logrus.Fields{
			"http.response.status_code":  resp.StatusCode,
			"http.response.body.content": string(b),
			"http.response.headers":      resp.Header,
		}).Debug("clearbank.client -> response")
	}

	if !slices.Contains(req.expectedStatuses, resp.StatusCode) {
		unexpectedResponse := UnexpectedResponse{Status: resp.StatusCode, Body: string(b)}

		var errResponse ErrResponse
		if err := json.Unmarshal(b, &errResponse); err != nil {
			return unexpectedResponse
		}

		if err := c.validate.Struct(errResponse); err != nil {
			return unexpectedResponse
		}

		return errResponse
	}

	if req.decodeTo != nil {
		if err := json.NewDecoder(resp.Body).Decode(req.decodeTo); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
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
