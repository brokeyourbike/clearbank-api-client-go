package clearbank_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/brokeyourbike/clearbank-api-client-go"
	"github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFetchAccountMandates_RequestErr(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	_, err := client.FetchAccountMandates(nil, uuid.New(), 0, 0) //lint:ignore SA1012 testing failure
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create request")
}

func TestCancelAccountMandate_RequestErr(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	err := client.CancelAccountMandate(nil, uuid.New(), uuid.New(), "0") //lint:ignore SA1012 testing failure
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create request")
}

func TestFetchVirtualAccountMandates_RequestErr(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	_, err := client.FetchVirtualAccountMandates(nil, uuid.New(), uuid.New(), 0, 0) //lint:ignore SA1012 testing failure
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create request")
}

func TestCancelVirtualAccountMandate_RequestErr(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	err := client.CancelVirtualAccountMandate(nil, uuid.New(), uuid.New(), uuid.New(), "0") //lint:ignore SA1012 testing failure
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create request")
}

func TestReturnBacsPayment_RequestErr(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	client := clearbank.NewClient("token", nil, clearbank.WithHTTPClient(mockHttpClient))

	err := client.ReturnBacsPayment(nil, uuid.New(), clearbank.ReturnBacsPaymentPayload{}) //lint:ignore SA1012 testing failure
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create request")
}

func TestReturnBacsPayment(t *testing.T) {
	mockHttpClient := clearbank.NewMockHttpClient(t)
	mockSigner := clearbank.NewMockSigner(t)
	client := clearbank.NewClient("token", mockSigner, clearbank.WithHTTPClient(mockHttpClient))

	ctx := clearbank.RequestIdContext(context.TODO(), "123")
	mockSigner.On("Sign", ctx, mock.Anything).Return([]byte("signed"), nil).Once()

	resp := &http.Response{StatusCode: http.StatusAccepted, Body: io.NopCloser(bytes.NewReader(nil))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	err := client.ReturnBacsPayment(ctx, uuid.New(), clearbank.ReturnBacsPaymentPayload{
		Returns: []clearbank.BacsReturn{
			{TransactionID: uuid.New(), ReasonCode: "0"},
		},
	})
	require.NoError(t, err)
}
