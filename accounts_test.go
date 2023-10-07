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

//go:embed testdata/account-success.json
var accountSuccess []byte

//go:embed testdata/accounts-success.json
var accountsSuccess []byte

func TestFetchAccount(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(accountSuccess))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	accountID := uuid.MustParse("a85002e3-0116-4b14-b7fa-427e60f4f6bc")

	got, err := client.FetchAccount(context.TODO(), accountID)
	require.NoError(t, err)
	assert.Equal(t, accountID, got.Account.ID)
	assert.Len(t, got.Account.Balances, 1)
	assert.Len(t, got.Account.Currency, 1)
}

func TestFetchAccounts(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(accountsSuccess))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.FetchAccounts(context.TODO(), 1, 100)
	require.NoError(t, err)

	assert.Len(t, got.Accounts, 1)
}
