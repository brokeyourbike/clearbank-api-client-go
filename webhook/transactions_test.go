package webhook_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/brokeyourbike/clearbank-api-client-go/webhook"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/TransactionSettled-with-supplementary.json
var transactionSettledWithSupplementary []byte

//go:embed testdata/TransactionSettled-no-supplementary.json
var transactionSettledNoSupplementary []byte

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
