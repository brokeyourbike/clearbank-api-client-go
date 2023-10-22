package clearbank

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
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

func TestUUID(t *testing.T) {
	raw := "018223b0-b499-43c6-be4f-518c916d9256"
	id := uuid.MustParse(raw)

	//lint:ignore S1025 testing how uuid is printed
	assert.Equal(t, raw, fmt.Sprintf("%s", id))
}
