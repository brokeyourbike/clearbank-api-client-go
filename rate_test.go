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

//go:embed testdata/marketrate-success.json
var marketrateSuccess []byte

//go:embed testdata/marketrate-not-available.txt
var marketrateNotAvailable []byte

//go:embed testdata/marketrate-outside-trading-hours.txt
var marketrateOutsideTradingHours []byte

//go:embed testdata/negotiate-success.json
var negotiateSuccess []byte

//go:embed testdata/bad-request-invalid.json
var badRequestInvalid []byte

func TestFetchMarketrate(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(marketrateSuccess))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.FetchMarketrate(context.TODO(), clearbank.MarketrateParams{FixedSide: clearbank.FixedSideBuy})
	require.NoError(t, err)

	assert.Equal(t, 1.209602, got.MarketRate)
	assert.Equal(t, "GBP/USD", got.Symbol)
	assert.False(t, clearbank.MarketrateNotAvailableError(err))
}

func TestFetchMarketrate_RequestErr(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	_, err := client.FetchMarketrate(nil, clearbank.MarketrateParams{FixedSide: clearbank.FixedSideBuy}) //lint:ignore SA1012 testing failure
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create request")
}

func TestFetchMarketrate_NotAvailable(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: 404, Body: io.NopCloser(bytes.NewReader(marketrateNotAvailable))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	_, err := client.FetchMarketrate(context.TODO(), clearbank.MarketrateParams{FixedSide: clearbank.FixedSideBuy})
	require.Error(t, err)
	require.True(t, clearbank.MarketrateNotAvailableError(err))
}

func TestFetchMarketrate_OutsideTradingHours(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: 400, Body: io.NopCloser(bytes.NewReader(marketrateOutsideTradingHours))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	_, err := client.FetchMarketrate(context.TODO(), clearbank.MarketrateParams{FixedSide: clearbank.FixedSideBuy})
	require.Error(t, err)
	require.True(t, clearbank.MarketrateNotAvailableError(err))
}

func TestNegotiate(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(negotiateSuccess))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.Negotiate(context.TODO())
	require.NoError(t, err)

	assert.Equal(t, "https://example.com", got.URL)
}

func TestNegotiate_RequestErr(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	_, err := client.Negotiate(nil) //lint:ignore SA1012 testing failure
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create request")
}

func TestNegotiate_ValidationFailed(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: 123, Body: io.NopCloser(bytes.NewReader(badRequestInvalid))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	_, err := client.Negotiate(context.TODO())
	require.Error(t, err)
	require.ErrorIs(t, err, clearbank.UnexpectedResponse{Status: 123, Body: "{}"}, "err response with no required fields is unexpected")
}
