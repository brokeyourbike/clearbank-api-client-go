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

//go:embed testdata/accounts-success.json
var accountsSuccess []byte

func TestFetchAccounts(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(accountsSuccess))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.FetchAccounts(context.TODO(), 1, 100)
	require.NoError(t, err)

	assert.Len(t, got.Accounts, 1)
}
