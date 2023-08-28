package clearbank

import (
	"context"
	"testing"

	logrustest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	httpClient := NewMockHttpClient(t)
	logger, _ := logrustest.NewNullLogger()

	client := NewClient("", nil, WithBaseURL("https://c.om"), WithHTTPClient(httpClient), WithLogger(logger))

	assert.Equal(t, "https://c.om", client.baseURL)
	assert.Same(t, httpClient, client.httpClient)
	assert.Same(t, logger, client.logger)
}

func TestRequestIdFromContext(t *testing.T) {
	assert.Equal(t, "", RequestIdFromContext(nil)) //lint:ignore SA1012 this is a test
	assert.Equal(t, "", RequestIdFromContext(context.TODO()))
	assert.Equal(t, "123", RequestIdFromContext(RequestIdContext(context.TODO(), "123")))
}
