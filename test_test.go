package clearbank_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/brokeyourbike/clearbank-api-client-go"
	"github.com/brokeyourbike/clearbank-api-client-go/signature"
	"github.com/sirupsen/logrus"
	logrustest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTest(t *testing.T) {
	mockSigner := signature.NewMockSigner(t)
	mockHttpClient := clearbank.NewMockHttpClient(t)

	logger, hook := logrustest.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	client := clearbank.NewClient("token", mockSigner, clearbank.WithHTTPClient(mockHttpClient), clearbank.WithLogger(logger))

	ctx := context.TODO()

	mockSigner.On("Sign", ctx, mock.Anything).Return([]byte("signed"), nil).Once()

	resp := &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(nil))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	assert.NoError(t, client.Test(ctx, "hello!"))

	require.Equal(t, 2, len(hook.Entries))
	require.Contains(t, hook.Entries[0].Data, "http.request.body.content")
	require.Contains(t, hook.Entries[0].Data, "http.request.headers.request_id")
	require.Contains(t, hook.Entries[1].Data, "http.response.status_code")
	require.Contains(t, hook.Entries[1].Data, "http.response.body.content")
	require.Contains(t, hook.Entries[1].Data, "http.response.headers")
}