package webhook_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/brokeyourbike/clearbank-api-client-go/webhook"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/BacsMandateCancelledV2.json
var bacsMandateCancelledV2 []byte

//go:embed testdata/BacsMandateInitiatedV2.json
var bacsMandateInitiatedV2 []byte

func TestBacsMandateCancelledV2Payload(t *testing.T) {
	var d1 webhook.BacsMandateCancelledV2Payload
	err := json.Unmarshal(bacsMandateCancelledV2, &d1)
	assert.NoError(t, err)
	assert.Equal(t, uuid.MustParse("8a1c544f-304b-d55a-f459-98d0bedaf19f"), d1.MandateID)
	assert.Equal(t, "MARY BROWN", d1.PayerInformation.AccountName)
}

func TestBacsMandateInitiatedV2Payload(t *testing.T) {
	var d1 webhook.BacsMandateInitiatedV2Payload
	err := json.Unmarshal(bacsMandateInitiatedV2, &d1)
	assert.NoError(t, err)
	assert.Equal(t, uuid.MustParse("086ab6ec-ab78-2432-12a1-3c273d81d961"), d1.MandateID)
	assert.Equal(t, "MARY BROWN", d1.PayerInformation.AccountName)
}
