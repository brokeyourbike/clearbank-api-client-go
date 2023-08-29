package clearbank_test

import (
	"bytes"
	"context"
	_ "embed"
	"io"
	"net/http"
	"testing"

	"github.com/brokeyourbike/clearbank-api-client-go"
	"github.com/brokeyourbike/clearbank-api-client-go/signature"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/marketrate-success.json
var marketrateSuccess []byte

//go:embed testdata/marketrate-fail.txt
var marketrateFail []byte

//go:embed testdata/negotiate-success.json
var negotiateSuccess []byte

func TestFetchMarketrate(t *testing.T) {
	mockSigner := signature.NewMockSigner(t)
	mockHttpClient := clearbank.NewMockHttpClient(t)

	client := clearbank.NewClient("token", mockSigner, clearbank.WithHTTPClient(mockHttpClient))

	ctx := context.TODO()
	mockSigner.On("Sign", ctx, mock.Anything).Return([]byte("signed"), nil).Once()

	resp := &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(marketrateSuccess))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.FetchMarketrate(ctx, clearbank.MarketrateParams{FixedSide: clearbank.FixedSideBuy})
	require.NoError(t, err)

	assert.Equal(t, 1.209602, got.MarketRate)
	assert.Equal(t, "GBP/USD", got.Symbol)
}

func TestFetchMarketrate_Fail(t *testing.T) {
	mockSigner := signature.NewMockSigner(t)
	mockHttpClient := clearbank.NewMockHttpClient(t)

	client := clearbank.NewClient("token", mockSigner, clearbank.WithHTTPClient(mockHttpClient))

	ctx := context.TODO()
	mockSigner.On("Sign", ctx, mock.Anything).Return([]byte("signed"), nil).Once()

	resp := &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(marketrateFail))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	_, err := client.FetchMarketrate(ctx, clearbank.MarketrateParams{FixedSide: clearbank.FixedSideBuy})
	require.Error(t, err)
}

func TestFetchNegotiate(t *testing.T) {
	mockSigner := signature.NewMockSigner(t)
	mockHttpClient := clearbank.NewMockHttpClient(t)

	client := clearbank.NewClient("token", mockSigner, clearbank.WithHTTPClient(mockHttpClient))

	ctx := context.TODO()
	mockSigner.On("Sign", ctx, mock.Anything).Return([]byte("signed"), nil).Once()

	resp := &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(negotiateSuccess))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.Negotiate(ctx)
	require.NoError(t, err)

	assert.Equal(t, "https://example.com", got.URL)
}
