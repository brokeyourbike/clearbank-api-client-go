package clearbank_test

import (
	"bytes"
	"context"
	_ "embed"
	"io"
	"net/http"
	"slices"
	"testing"

	"github.com/brokeyourbike/clearbank-api-client-go"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/fx-request-fx-quote-success.json
var requestFxQuoteSuccess []byte

func TestInitiateFxOrder(t *testing.T) {
	mockSigner := clearbank.NewMockSigner(t)
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", mockSigner, clearbank.WithHTTPClient(mockHttpClient))

	ctx := clearbank.RequestIdContext(context.TODO(), "123")
	mockSigner.On("Sign", ctx, mock.Anything).Return([]byte("signed"), nil).Once()

	resp := &http.Response{StatusCode: http.StatusAccepted, Body: io.NopCloser(bytes.NewReader(nil))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	assert.NoError(t, client.InitiateFxOrder(ctx, clearbank.FXPayload{}))
}

func TestInitiateFxOrder_RequestErr(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	err := client.InitiateFxOrder(nil, clearbank.FXPayload{}) //lint:ignore SA1012 testing failure
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create request")
}

func TestExecuteFxQuote(t *testing.T) {
	mockSigner := clearbank.NewMockSigner(t)
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", mockSigner, clearbank.WithHTTPClient(mockHttpClient))

	ctx := clearbank.RequestIdContext(context.TODO(), "123")
	mockSigner.On("Sign", ctx, mock.Anything).Return([]byte("signed"), nil).Once()

	resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(nil))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	assert.NoError(t, client.ExecuteFxQuote(ctx, clearbank.FXPayload{}))
}

func TestExecuteFxQuote_RequestIdFromCtx(t *testing.T) {
	mockSigner := clearbank.NewMockSigner(t)
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", mockSigner, clearbank.WithHTTPClient(mockHttpClient))

	ctx := clearbank.RequestIdContext(context.TODO(), "123")
	mockSigner.On("Sign", ctx, mock.Anything).Return([]byte("signed"), nil).Once()

	resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(nil))}
	mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool { return slices.Contains(req.Header["X-Request-Id"], "123") })).Return(resp, nil).Once()

	assert.NoError(t, client.ExecuteFxQuote(ctx, clearbank.FXPayload{}))
}

func TestExecuteFxQuote_RequestIdFromFunc(t *testing.T) {
	mockSigner := clearbank.NewMockSigner(t)
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", mockSigner,
		clearbank.WithHTTPClient(mockHttpClient),
		clearbank.WithRequestIDGenerator(func() string { return "456" }),
	)

	ctx := context.TODO()
	mockSigner.On("Sign", ctx, mock.Anything).Return([]byte("signed"), nil).Once()

	resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(nil))}
	mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool { return slices.Contains(req.Header["X-Request-Id"], "456") })).Return(resp, nil).Once()

	assert.NoError(t, client.ExecuteFxQuote(ctx, clearbank.FXPayload{}))
}

func TestExecuteFxQuote_RequestIdFromCtxAndFunc(t *testing.T) {
	mockSigner := clearbank.NewMockSigner(t)
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", mockSigner,
		clearbank.WithHTTPClient(mockHttpClient),
		clearbank.WithRequestIDGenerator(func() string { return "456" }),
	)

	ctx := clearbank.RequestIdContext(context.TODO(), "123")
	mockSigner.On("Sign", ctx, mock.Anything).Return([]byte("signed"), nil).Once()

	resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(nil))}
	// ctx takex priority and overrides the func
	mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool { return slices.Contains(req.Header["X-Request-Id"], "123") })).Return(resp, nil).Once()

	assert.NoError(t, client.ExecuteFxQuote(ctx, clearbank.FXPayload{}))
}

func TestRequestFxQuote(t *testing.T) {
	mockSigner := clearbank.NewMockSigner(t)
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", mockSigner, clearbank.WithHTTPClient(mockHttpClient))

	ctx := clearbank.RequestIdContext(context.TODO(), "123")
	mockSigner.On("Sign", ctx, mock.Anything).Return([]byte("signed"), nil).Once()

	resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(requestFxQuoteSuccess))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.RequestFxQuote(ctx, clearbank.FXQuotePayload{})
	require.NoError(t, err)
	assert.Equal(t, "32fd7f3f-2abe-4f83-b15d-1285a07f6053", got.QuoteID.String())
}

func TestFXQuoteResponse(t *testing.T) {
	r1 := clearbank.FXQuoteResponse{SellCurrency: "EUR", BuyCurrency: "GBP", CurrencyPair: "EUR/GBP", ExchangeRate: 1.23}
	require.Equal(t, 1.23, r1.GetRate())

	r2 := clearbank.FXQuoteResponse{SellCurrency: "EUR", BuyCurrency: "GBP", CurrencyPair: "GBP/EUR", ExchangeRate: 1.23}
	require.Equal(t, 1/1.23, r2.GetRate())
}
