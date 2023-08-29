package clearbank

import (
	"context"
	"fmt"
	"net/http"
)

type StatementClient interface {
	RequestStatement(context.Context, StatementPayload) error
	RequestStatementFor(context.Context, string, StatementPayload) error
}

type StatementPayload struct {
	Year     int    `json:"year"`
	Month    int    `json:"month"`
	Format   string `json:"format"`
	Currency string `json:"currency"`
}

func (c *client) RequestStatement(ctx context.Context, payload StatementPayload) error {
	req, err := c.newRequest(ctx, http.MethodPost, "/mccy/v1/StatementRequests", payload)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusAccepted)
	return c.do(ctx, req)
}

func (c *client) RequestStatementFor(ctx context.Context, iban string, payload StatementPayload) error {
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/mccy/v1/StatementRequests/account/%s", iban), payload)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusAccepted)
	return c.do(ctx, req)
}
