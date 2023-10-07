package stream_test

import (
	"testing"

	"github.com/brokeyourbike/clearbank-api-client-go/stream"
	"github.com/stretchr/testify/assert"
)

func TestSubscriptionRequest(t *testing.T) {
	req := stream.NewSubscriptionRequest("A/B")
	assert.Equal(t, 1, req.Version)
	assert.Equal(t, "SubscriptionRequest", req.Type)
	assert.Equal(t, "A/B", req.Payload.CurrencyPair)
}
