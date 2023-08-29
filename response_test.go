package clearbank_test

import (
	"testing"

	"github.com/brokeyourbike/clearbank-api-client-go"
	"github.com/stretchr/testify/assert"
)

func TestUnexpectedResponse(t *testing.T) {
	err := clearbank.UnexpectedResponse{Status: 123}
	assert.Contains(t, err.Error(), "status: 123")
}

func TestErrResponse(t *testing.T) {
	err := clearbank.ErrResponse{Status: 123, Type: "t", Title: "hi"}
	assert.Contains(t, err.Error(), "status: 123")
	assert.Contains(t, err.Error(), "type: t")
	assert.Contains(t, err.Error(), "title: hi")
}
