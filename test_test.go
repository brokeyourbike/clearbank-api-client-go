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
	"github.com/stretchr/testify/mock"
)

func TestTest(t *testing.T) {
	mockSigner := signature.NewMockSigner(t)
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", mockSigner, clearbank.WithHTTPClient(mockHttpClient))

	ctx := context.TODO()

	mockSigner.On("Sign", ctx, mock.Anything).Return([]byte("signed"), nil).Once()

	resp := &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(nil))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	assert.NoError(t, client.Test(ctx, "hello!"))
}
