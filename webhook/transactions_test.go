package webhook_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/brokeyourbike/clearbank-api-client-go/webhook"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/webhook-TransactionSettled-with-supplementary.json
var webhookTransactionSettledWithSupplementary []byte

//go:embed testdata/webhook-TransactionSettled-no-supplementary.json
var webhookTransactionSettledNoSupplementary []byte

func TestWebhookTransactionSettledPayload(t *testing.T) {
	var d1 webhook.WebhookTransactionSettledPayload
	err := json.Unmarshal(webhookTransactionSettledWithSupplementary, &d1)
	assert.NoError(t, err)
	assert.Len(t, d1.SupplementaryData, 1)

	var d2 webhook.WebhookTransactionSettledPayload
	err = json.Unmarshal(webhookTransactionSettledNoSupplementary, &d2)
	assert.NoError(t, err)
	assert.Len(t, d2.SupplementaryData, 0)
}
