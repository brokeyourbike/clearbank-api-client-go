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

//go:embed testdata/account-create-success.json
var accountCreateSuccess []byte

//go:embed testdata/account-fetch-virtual-account.json
var accountFetchVirtualAccountSuccess []byte

func TestFetchAccount(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(accountSuccess))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	accountID := uuid.MustParse("a85002e3-0116-4b14-b7fa-427e60f4f6bc")

	got, err := client.FetchAccount(context.TODO(), accountID)
	require.NoError(t, err)
	assert.Equal(t, accountID, got.Account.ID)
	assert.Len(t, got.Account.Balances, 1)
	assert.Len(t, got.Account.Currency, 1)
}

func TestFetchAccount_RequestErr(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	_, err := client.FetchAccount(nil, uuid.New()) //lint:ignore SA1012 testing failure
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create request")
}

func TestFetchAccounts(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(accountsSuccess))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.FetchAccounts(context.TODO(), 1, 100)
	require.NoError(t, err)

	assert.Len(t, got.Accounts, 1)
}

func TestFetchAccounts_RequestErr(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	_, err := client.FetchAccounts(nil, 1, 100) //lint:ignore SA1012 testing failure
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create request")
}

func TestCreateAccount(t *testing.T) {
	mockSigner := clearbank.NewMockSigner(t)
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", mockSigner, clearbank.WithHTTPClient(mockHttpClient))

	ctx := clearbank.RequestIdContext(context.TODO(), "123")
	mockSigner.On("Sign", ctx, mock.Anything).Return([]byte("signed"), nil).Once()

	resp := &http.Response{StatusCode: http.StatusCreated, Body: io.NopCloser(bytes.NewReader(accountCreateSuccess))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.CreateAccount(ctx, clearbank.CreateAccountPayload{})
	require.NoError(t, err)
	assert.Equal(t, uuid.MustParse("a85002e3-0116-4b14-b7fa-427e60f4f6bc"), got.Account.ID)
}

func TestCreateAccount_RequestErr(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	_, err := client.CreateAccount(nil, clearbank.CreateAccountPayload{}) //lint:ignore SA1012 testing failure
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create request")
}

func TestFetchVirtualAccount(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(accountFetchVirtualAccountSuccess))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	accountID := uuid.MustParse("9cd1ae28-1db1-447c-8775-32e355234e9d")
	virtualAccountID := uuid.MustParse("8c3748d3-48c1-4e20-a092-cf09d3faf58c")

	got, err := client.FetchVirtualAccount(context.TODO(), accountID, virtualAccountID)
	require.NoError(t, err)
	assert.Equal(t, virtualAccountID, got.Account.ID)
	assert.Equal(t, clearbank.VirtualAccountStatusDisabled, got.Account.Status)
}
