package clearbank_test

import (
	"testing"

	"github.com/brokeyourbike/clearbank-api-client-go"
	"github.com/google/uuid"
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
