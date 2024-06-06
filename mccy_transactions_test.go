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

//go:embed testdata/mccy-transaction-with-virtual-account.json
var mccyTransactionWithVirtualAccountMsg []byte

func TestInitiateInternalTransaction_RequestErr(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	err := client.InitiateInternalTransaction(nil, clearbank.CreateInternalTransactionPayload{}) //lint:ignore SA1012 testing failure
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create request")
}

func TestInitiateMCCYTransactions_RequestErr(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	err := client.InitiateMCCYTransactions(nil, clearbank.CreateMCCYTransactionsPayload{}) //lint:ignore SA1012 testing failure
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create request")
}

func TestInitiateMCCYInboundPayment_RequestErr(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	_, err := client.InitiateMCCYInboundPayment(nil, "currency", clearbank.CreateMCCYInboundPaymentPayload{}) //lint:ignore SA1012 testing failure
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create request")
}

func TestFetchMCCYTransaction_WithVirtualAccount(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)

	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))
	ctx := clearbank.RequestIdContext(context.TODO(), "123")

	resp := &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(mccyTransactionWithVirtualAccountMsg))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	trx, err := client.FetchMCCYTransaction(ctx, uuid.New())
	require.NoError(t, err)

	assert.Equal(t, "cbe94f68-ffa4-452e-99af-024b13fc56c8", trx.TransactionID.String())
	assert.Equal(t, "0f129da0-b072-4d8a-8cf4-3b7b879a7557", trx.AccountID.String())
	assert.Equal(t, "8e599dc8-cd35-49e4-81cc-b6ddd6ef7d16", trx.VirtualAccountID.UUID.String())
}

func TestFetchMCCYTransaction_RequestErr(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	_, err := client.FetchMCCYTransaction(nil, uuid.New()) //lint:ignore SA1012 testing failure
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create request")
}

func TestFetchMCCYTransactionsForAccount_RequestErr(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	_, err := client.FetchMCCYTransactionsForAccount(nil, uuid.New(), "USD", clearbank.FetchTransactionsParams{}) //lint:ignore SA1012 testing failure
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create request")
}

func TestFetchMCCYTransactionsForVirtualAccount_RequestErr(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	_, err := client.FetchMCCYTransactionsForVirtualAccount(nil, uuid.New(), "USD", clearbank.FetchTransactionsParams{}) //lint:ignore SA1012 testing failure
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create request")
}
