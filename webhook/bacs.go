package webhook

import "github.com/google/uuid"

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
