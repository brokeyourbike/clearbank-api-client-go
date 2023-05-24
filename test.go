package clearbank

import (
	"context"
)

type TestService struct {
	c *Client
}

func (s *TestService) Test(ctx context.Context, message string) error {
	// req, err := s.c.NewRequest(ctx, http.MethodPost, "/v1/Test", TestPayload{Data: message})
	// if err != nil {
	// 	return err
	// }

	return nil
}

// func (c *Client) Test(ctx context.Context, message string) error {
// 	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseUrl+"/mccy/v1/StatementRequests/account/%s", nil)
// 	if err != nil {
// 		return fmt.Errorf("failed to create request: %w", err)
// 	}

// 	if req.Method == http.MethodPost || req.Method == http.MethodPatch {
// 		body, err := json.Marshal(TestPayload{Data: message})
// 		if err != nil {
// 			return fmt.Errorf("failed to marshal payload: %w", err)
// 		}

// 		signature, err := c.signer.Sign(ctx, body)
// 		if err != nil {
// 			return fmt.Errorf("failed to sign payload: %w", err)
// 		}

// 		req.Body = io.NopCloser(bytes.NewReader(body))
// 		req.Header.Set("DigitalSignature", string(signature))
// 	}

// 	req.Header.Set("Content-Type", "application/json")
// 	// req.Header.Set("X-Request-Id", uuid.NewString()) // TODO: we can add a function to, or add a context key
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

// 	resp, err := c.httpClient.Do(req)
// 	if err != nil {
// 		return fmt.Errorf("failed to send request: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	b, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return fmt.Errorf("failed to read response body: %w", err)
// 	}

// 	resp.Body = io.NopCloser(bytes.NewBuffer(b))

// 	if resp.StatusCode != expectedStatus {
// 		// unexpectedResponse := UnexpectedResponse{Status: resp.StatusCode, Body: string(b)}

// 		var errResponse ErrResponse
// 		if err := json.Unmarshal(b, &errResponse); err != nil {
// 			return resp, unexpectedResponse
// 		}

// 		if errResponse.Status == 0 {
// 			return resp, unexpectedResponse
// 		}
// 	}

// 	return nil
// }
