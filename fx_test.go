package clearbank_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/brokeyourbike/clearbank-api-client-go"
	"github.com/brokeyourbike/clearbank-api-client-go/signature"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestInitiateFxOrder(t *testing.T) {
	mockSigner := signature.NewMockSigner(t)
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", mockSigner, clearbank.WithHTTPClient(mockHttpClient))

	ctx := clearbank.RequestIdContext(context.TODO(), "123")
	mockSigner.On("Sign", ctx, mock.Anything).Return([]byte("signed"), nil).Once()

	resp := &http.Response{StatusCode: http.StatusAccepted, Body: io.NopCloser(bytes.NewReader(nil))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	assert.NoError(t, client.InitiateFxOrder(ctx, clearbank.FXPayload{}))
}
