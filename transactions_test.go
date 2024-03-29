package clearbank_test

import (
	"bytes"
	"context"
	_ "embed"
	"io"
	"net/http"
	"testing"

	"github.com/brokeyourbike/clearbank-api-client-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/transactions-fetch-one-for-account.json
var transactionsFetchOneForAccount []byte

func TestFetchTransactionForAccount(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(transactionsFetchOneForAccount))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.FetchTransactionForAccount(context.TODO(), uuid.New(), uuid.New())
	assert.NoError(t, err)
	assert.Equal(t, "6ea6dd13-eaad-4e1a-990c-6533a177dbb4", got.Transaction.TransactionID.String())
}

func TestFetchTransactionForAccount_RequestErr(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	_, err := client.FetchTransactionForAccount(nil, uuid.New(), uuid.New()) //lint:ignore SA1012 testing failure
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create request")
}

func TestFetchTransactionForVirtualAccount(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(transactionsFetchOneForAccount))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.FetchTransactionForVirtualAccount(context.TODO(), uuid.New(), uuid.New(), uuid.New())
	assert.NoError(t, err)
	assert.Equal(t, "6ea6dd13-eaad-4e1a-990c-6533a177dbb4", got.Transaction.TransactionID.String())
}

func TestFetchTransactionForVirtualAccount_RequestErr(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	_, err := client.FetchTransactionForVirtualAccount(nil, uuid.New(), uuid.New(), uuid.New()) //lint:ignore SA1012 testing failure
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create request")
}
