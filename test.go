package clearbank

import (
	"context"
	"fmt"
	"net/http"
)

type TestClient interface {
	Test(context.Context, string) error
}

type TestPayload struct {
	Data string `json:"data"`
}

func (c *client) Test(ctx context.Context, message string) error {
	req, err := c.NewRequest(ctx, http.MethodPost, "/v1/Test", TestPayload{Data: message})
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if _, err := c.Do(ctx, req, http.StatusOK); err != nil {
		return err
	}

	return nil
}
