package webhook_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/brokeyourbike/clearbank-api-client-go/webhook"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/TransactionSettled-with-supplementary.json
var transactionSettledWithSupplementary []byte

//go:embed testdata/TransactionSettled-no-supplementary.json
var transactionSettledNoSupplementary []byte

//go:embed testdata/InboundHeldTransaction.json
var inboundHeldTransaction []byte

//go:embed testdata/OutboundHeldTransaction.json
var outboundHeldTransaction []byte

//go:embed testdata/OutboundHeldTransaction-no-time.json
var outboundHeldTransactionNoTime []byte

func TestTransactionSettledPayload(t *testing.T) {
	var d1 webhook.TransactionSettledPayload
	err := json.Unmarshal(transactionSettledWithSupplementary, &d1)
	assert.NoError(t, err)
	assert.Len(t, d1.SupplementaryData, 1)

	var d2 webhook.TransactionSettledPayload
	err = json.Unmarshal(transactionSettledNoSupplementary, &d2)
	assert.NoError(t, err)
	assert.Len(t, d2.SupplementaryData, 0)
}

func TestHeldTransaction(t *testing.T) {
	var d1 webhook.InboundHeldTransactionPayload
	err := json.Unmarshal(inboundHeldTransaction, &d1)
	require.NoError(t, err)
	assert.Equal(t, "GB00CUBK11223312345678", d1.Account.IBAN)

	var d2 webhook.OutboundHeldTransactionPayload
	err = json.Unmarshal(outboundHeldTransaction, &d2)
	require.NoError(t, err)
	assert.Equal(t, "GB00CUBK11223312345678", d2.Account.IBAN)

	var d3 webhook.OutboundHeldTransactionPayload
	err = json.Unmarshal(outboundHeldTransactionNoTime, &d3)
	require.NoError(t, err)
	assert.Equal(t, "GB00CUBK11223312345678", d3.Account.IBAN)
}
