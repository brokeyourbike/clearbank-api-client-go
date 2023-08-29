package clearbank_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/brokeyourbike/clearbank-api-client-go"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/bad-request-validation.json
var badRequestValidation []byte

//go:embed testdata/bad-request-validation-no-errors.json
var badRequestValidationNoErrors []byte

func TestUnexpectedResponse(t *testing.T) {
	resp := clearbank.UnexpectedResponse{Status: 500, Body: "I am an error."}
	assert.Equal(t, resp.Error(), "Unexpected response from API. status: 500 body: I am an error.")
}

func TestErrResponse(t *testing.T) {
	resp := clearbank.ErrResponse{Status: 500, Type: "abcd", Title: "I am an error."}
	assert.Equal(t, resp.Error(), "Error during API call. status: 500 type: abcd title: I am an error.")
}

func TestErrResponse_Errors(t *testing.T) {
	var resp clearbank.ErrResponse
	err := json.Unmarshal(badRequestValidation, &resp)
	assert.NoError(t, err)

	assert.Len(t, resp.Errors, 2)
	assert.Equal(t, resp.Errors["DebitAccountIban"][0], "Debit Account IBAN must be populated and should be in a valid IBAN format")
	assert.Equal(t, resp.Errors["CreditAccountIban"][0], "Credit Account IBAN must be populated and should be in a valid IBAN format")
}

func TestErrResponse_NoErrors(t *testing.T) {
	var resp clearbank.ErrResponse
	err := json.Unmarshal(badRequestValidationNoErrors, &resp)
	assert.NoError(t, err)

	assert.Len(t, resp.Errors, 0)
}
