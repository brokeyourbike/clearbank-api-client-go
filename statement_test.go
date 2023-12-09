package clearbank_test

import (
	"bytes"
	"context"
	_ "embed"
	"io"
	"net/http"
	"testing"

	"github.com/brokeyourbike/clearbank-api-client-go"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRequestStatement(t *testing.T) {
	mockSigner := clearbank.NewMockSigner(t)
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", mockSigner, clearbank.WithHTTPClient(mockHttpClient))

	ctx := clearbank.RequestIdContext(context.TODO(), "123")
	mockSigner.On("Sign", ctx, mock.Anything).Return([]byte("signed"), nil).Once()

	resp := &http.Response{StatusCode: http.StatusAccepted, Body: io.NopCloser(bytes.NewReader(nil))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	assert.NoError(t, client.RequestStatement(ctx, clearbank.StatementPayload{}))
}

func TestRequestStatement_RequestErr(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	err := client.RequestStatement(nil, clearbank.StatementPayload{}) //lint:ignore SA1012 testing failure
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create request")
}

func TestRequestStatementFor(t *testing.T) {
	mockSigner := clearbank.NewMockSigner(t)
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", mockSigner, clearbank.WithHTTPClient(mockHttpClient))

	ctx := clearbank.RequestIdContext(context.TODO(), "123")
	mockSigner.On("Sign", ctx, mock.Anything).Return([]byte("signed"), nil).Once()

	resp := &http.Response{StatusCode: http.StatusAccepted, Body: io.NopCloser(bytes.NewReader(nil))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	assert.NoError(t, client.RequestStatementFor(ctx, "IBAN12345", clearbank.StatementPayload{}))
}

func TestRequestStatementFor_RequestErr(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	err := client.RequestStatementFor(nil, "IBAN12345", clearbank.StatementPayload{}) //lint:ignore SA1012 testing failure
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create request")
}
