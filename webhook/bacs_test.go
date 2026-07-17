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

//go:embed testdata/BacsDirectDebitInboundPaymentCreatedV2.json
var bacsDirectDebitInboundPaymentCreatedV2 []byte

//go:embed testdata/BacsDirectDebitReturnCreatedV4.json
var bacsDirectDebitReturnCreatedV4 []byte

func TestBacsDirectDebitReturnCreatedV4Payload(t *testing.T) {
	var d1 webhook.BacsDirectDebitReturnCreatedV4Payload
	err := json.Unmarshal(bacsDirectDebitReturnCreatedV4, &d1)
	assert.NoError(t, err)
	assert.Equal(t, uuid.MustParse("14a216fa-a70b-a211-3f31-47c2d992c86c"), d1.BacsTransactionId)
	assert.Equal(t, "MARY BROWN", d1.PayerInformation.AccountName)
	assert.Equal(t, "Q-MOBILE", d1.OriginatorInformation.AccountName)
	assert.Equal(t, "B", d1.ReasonCode)
	assert.Equal(t, "2a51971a3576b2a5960ff65202e1539d", d1.Reference)
	assert.Equal(t, 25.0, d1.Amount)
	assert.Equal(t, "17", d1.TransactionCode)
	assert.Equal(t, "Bacs", d1.Source)
	assert.Equal(t, "2019-04-01 00:00:00 +0000 UTC", d1.ProcessingDate.String())
	assert.Equal(t, "030201", d1.ServiceUserNumber)
	assert.Equal(t, uuid.MustParse("d61b2513-f586-43da-8b72-4f6ba2e28fea"), d1.OriginalBacsTransactionId)
	assert.Equal(t, "012345", d1.OriginalBacsTransactionSun)
}

func TestBacsDirectDebitInboundPaymentCreatedV2Payload(t *testing.T) {
	var d1 webhook.BacsDirectDebitInboundPaymentCreatedV2Payload
	err := json.Unmarshal(bacsDirectDebitInboundPaymentCreatedV2, &d1)
	assert.NoError(t, err)
	assert.Equal(t, uuid.MustParse("d9c36960-53dc-fb00-8b4b-d35bdcd8ac03"), d1.BacsTransactionId)
	assert.Equal(t, "JOHN DOE", d1.PayerInformation.AccountName)
	assert.Equal(t, "UK TELECOMMS", d1.OriginatorInformation.AccountName)
	assert.Equal(t, 13.75, d1.Amount)
	assert.Equal(t, "JODTR54PXZNOV2018", d1.Reference)
	assert.Equal(t, "010203", d1.ServiceUserNumber)
	assert.Equal(t, "01", d1.TransactionCode)
	assert.Equal(t, "Bacs", d1.Source)
	assert.Equal(t, "2019-01-01 00:00:00 +0000 UTC", d1.ProcessingDate.String())
}

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
