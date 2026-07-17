package webhook

import (
	"github.com/brokeyourbike/clearbank-api-client-go"
	"github.com/google/uuid"
)

// BacsDirectDebitInboundPaymentCreatedV2
// This webhook confirms that a Bacs Direct Debit inbound payment has been received.
type BacsDirectDebitInboundPaymentCreatedV2Payload struct {
	PayerInformation      *BacsActor     `json:"PayerInformation"`
	OriginatorInformation *BacsActor     `json:"OriginatorInformation"`
	Amount                float64        `json:"Amount"`
	Reference             string         `json:"Reference"`
	ProcessingDate        clearbank.Time `json:"ProcessingDate"`
	ServiceUserNumber     string         `json:"ServiceUserNumber"`
	TransactionCode       string         `json:"TransactionCode"`
	Source                string         `json:"Source"`
	BacsTransactionId     uuid.UUID      `json:"BacsTransactionId"`
}

type ARUDDReasonCode string

const (
	ARUDDReferToPayer          ARUDDReasonCode = "0" // Refer to payer
	ARUDDInstructionCancelled  ARUDDReasonCode = "1" // Instruction cancelled
	ARUDDPayerDeceased         ARUDDReasonCode = "2" // Payer deceased
	ARUDDAccountTransferred    ARUDDReasonCode = "3" // Account transferred
	ARUDDAdvanceNoticeDisputed ARUDDReasonCode = "4" // Advance notice disputed
	ARUDDNoAccount             ARUDDReasonCode = "5" // No account (OR wrong account type)
	ARUDDNoInstruction         ARUDDReasonCode = "6" // No Instruction held
	ARUDDAmountDiffers         ARUDDReasonCode = "7" // Amount differs
	ARUDDAmountNotYetDue       ARUDDReasonCode = "8" // Amount not yet due
	ARUDDPresentationOverdue   ARUDDReasonCode = "9" // Presentation overdue
	ARUDDServiceUserDiffers    ARUDDReasonCode = "A" // Service user differs
	ARUDDAccountClosed         ARUDDReasonCode = "B" // Account closed
)

// BacsDirectDebitReturnCreatedV4
// This webhook confirms a Bacs Direct Debit return has been created.
type BacsDirectDebitReturnCreatedV4Payload struct {
	PayerInformation           *BacsActor      `json:"PayerInformation"`
	OriginatorInformation      *BacsActor      `json:"OriginatorInformation"`
	ReasonCode                 ARUDDReasonCode `json:"ReasonCode"`
	Reference                  string          `json:"Reference"`
	Amount                     float64         `json:"Amount"`
	TransactionCode            string          `json:"TransactionCode"`
	Source                     string          `json:"Source"`
	BacsTransactionId          uuid.UUID       `json:"BacsTransactionId"`
	ProcessingDate             clearbank.Time  `json:"ProcessingDate"`
	ServiceUserNumber          string          `json:"ServiceUserNumber"`
	OriginalBacsTransactionId  uuid.UUID       `json:"OriginalBacsTransactionId"`
	OriginalBacsTransactionSun string          `json:"OriginalBacsTransactionSun"`
}

// BacsMandateInitiatedV2
// This webhook confirms that a Bacs Direct Debit Instruction has been set up successfully.
type BacsMandateInitiatedV2Payload struct {
	MandateID             uuid.UUID  `json:"MandateId" validate:"required"`
	Reference             string     `json:"Reference"`
	Source                string     `json:"Source"`
	OriginatorInformation *BacsActor `json:"OriginatorInformation"`
	PayerInformation      *BacsActor `json:"PayerInformation"`
}

// BacsMandateCancelledV2
// This webhook confirms that a Bacs Direct Debit Instruction has been cancelled.
type BacsMandateCancelledV2Payload struct {
	MandateID             uuid.UUID  `json:"MandateId" validate:"required"`
	Reference             string     `json:"Reference"`
	Source                string     `json:"Source"`
	OriginatorInformation *BacsActor `json:"OriginatorInformation"`
	PayerInformation      *BacsActor `json:"PayerInformation"`
}

type BacsActor struct {
	AccountName   string `json:"AccountName"`
	SortCode      string `json:"SortCode"`
	AccountNumber string `json:"AccountNumber"`
}

// BacsMandateCancellationFailed
// This webhook confirms that a Bacs Direct Debit Instruction cancellation has failed.
type BacsMandateCancellationFailedPayload struct {
	ReasonCode            string     `json:"ReasonCode"`
	TransactionCode       string     `json:"TransactionCode"`
	Reference             string     `json:"Reference"`
	Source                string     `json:"Source"`
	OriginatorInformation *BacsActor `json:"OriginatorInformation"`
	PayerInformation      *BacsActor `json:"PayerInformation"`
}
