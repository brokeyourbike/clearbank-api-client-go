package clearbank_test

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
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

//go:embed testdata/marketrate-fail.txt
var marketrateFail []byte

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
}

func TestFetchMarketrate_Fail(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)

	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(marketrateFail))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	_, err := client.FetchMarketrate(context.TODO(), clearbank.MarketrateParams{FixedSide: clearbank.FixedSideBuy})
	require.Error(t, err)
}

func TestFetchMarketrate_FailedHttpRequest(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{}, errors.New("cannot do")).Once()

	_, err := client.FetchMarketrate(context.TODO(), clearbank.MarketrateParams{FixedSide: clearbank.FixedSideBuy})
	assert.Error(t, err)
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

func TestNegotiate_ValidationFailed(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: 123, Body: io.NopCloser(bytes.NewReader(badRequestInvalid))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	_, err := client.Negotiate(context.TODO())
	require.Error(t, err)
	require.ErrorIs(t, err, clearbank.UnexpectedResponse{Status: 123, Body: "{}"}, "err response with no required fields is unexpected")
}

func TestNegotiate_FailedHttpRequest(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{}, errors.New("cannot do")).Once()

	_, err := client.Negotiate(context.TODO())
	assert.Error(t, err)
}
