package webhook

import (
	"github.com/brokeyourbike/clearbank-api-client-go"
	"github.com/google/uuid"
)

type AdditionalProperty struct {
	Key   string `json:"Key" validate:"required"`
	Value string `json:"Value" validate:"required"`
}

// WebhookMCTransactionCreatedPayload
// This webhook confirms that a multicurrency transaction has been created
type WebhookMCTransactionCreatedPayload struct {
	TransactionID                      uuid.UUID            `json:"TransactionId" validate:"required"`
	BatchID                            uuid.UUID            `json:"BatchId" validate:"required"`
	EndToEndID                         string               `json:"EndToEndId" validate:"required"`
	SchemeEndToEndID                   string               `json:"SchemeEndToEndId"`
	AccountID                          uuid.UUID            `json:"AccountId" validate:"required"`
	Reference                          string               `json:"Reference" validate:"required"`
	UltimateCreditorName               string               `json:"UltimateCreditorName" validate:"required"`
	UltimateCreditorAccountIdentifiers []MCCYIdentifier     `json:"UltimateCreditorAccountIdentifiers" validate:"required"`
	UltimateDebtorName                 string               `json:"UltimateDebtorName" validate:"required"`
	UltimateDebtorAccountIdentifiers   []MCCYIdentifier     `json:"UltimateDebtorAccountIdentifiers" validate:"required"`
	InstructedAmount                   float64              `json:"InstructedAmount" validate:"required"`
	InstructedCurrency                 string               `json:"InstructedCurrency" validate:"required"`
	SchemePaymentMethod                string               `json:"SchemePaymentMethod" validate:"required"`
	TimestampCreated                   clearbank.Time       `json:"TimestampCreated" validate:"required"`
	AdditionalProperties               []AdditionalProperty `json:"AdditionalProperties" validate:"required"`
}

var _ Transaction = (*WebhookMCTransactionSettledPayload)(nil)

// WebhookMCTransactionSettledPayload
// This webhook confirms that a multicurrency transaction has been settled
type WebhookMCTransactionSettledPayload struct {
	TransactionID                      uuid.UUID            `json:"TransactionId" validate:"required"`
	BatchID                            uuid.UUID            `json:"BatchId" validate:"required"`
	EndToEndID                         string               `json:"EndToEndId" validate:"required"`
	SchemeEndToEndID                   string               `json:"SchemeEndToEndId"`
	AccountID                          uuid.UUID            `json:"AccountId" validate:"required"`
	Reference                          string               `json:"Reference" validate:"required"`
	UltimateCreditorName               string               `json:"UltimateCreditorName" validate:"required"`
	UltimateCreditorAccountIdentifiers []MCCYIdentifier     `json:"UltimateCreditorAccountIdentifiers" validate:"required"`
	UltimateDebtorName                 string               `json:"UltimateDebtorName" validate:"required"`
	UltimateDebtorAccountIdentifiers   []MCCYIdentifier     `json:"UltimateDebtorAccountIdentifiers" validate:"required"`
	InstructedAmount                   float64              `json:"InstructedAmount" validate:"required"`
	InstructedCurrency                 string               `json:"InstructedCurrency" validate:"required"`
	SchemePaymentMethod                string               `json:"SchemePaymentMethod" validate:"required"`
	TimestampCreated                   clearbank.Time       `json:"TimestampCreated" validate:"required"`
	TimestampSettled                   clearbank.Time       `json:"TimestampSettled" validate:"required"`
	TimestampSubmitted                 clearbank.Time       `json:"TimestampSubmitted" validate:"required"`
	AdditionalProperties               []AdditionalProperty `json:"AdditionalProperties" validate:"required"`
}

func (w WebhookMCTransactionSettledPayload) GetID() uuid.UUID {
	return w.TransactionID
}

func (w WebhookMCTransactionSettledPayload) GetEndToEndID() string {
	return w.EndToEndID
}

func (w WebhookMCTransactionSettledPayload) GetCurrency() string {
	return w.InstructedCurrency
}

func (w WebhookMCTransactionSettledPayload) GetAmount() float64 {
	return w.InstructedAmount
}

func (w WebhookMCTransactionSettledPayload) IsReturn() bool {
	return false
}

func (w WebhookMCTransactionSettledPayload) GetReference() string {
	return w.Reference
}

func (w WebhookMCTransactionSettledPayload) GetAccountIdentifier() string {
	if len(w.UltimateCreditorAccountIdentifiers) == 0 {
		return ""
	}
	if w.InstructedAmount > 0 {
		return w.UltimateCreditorAccountIdentifiers[0].Identifier
	}
	return w.UltimateDebtorAccountIdentifiers[0].Identifier
}

func (w WebhookMCTransactionSettledPayload) GetAccountOwner() string {
	if w.InstructedAmount > 0 {
		return w.UltimateCreditorName
	}
	return w.UltimateDebtorName
}

func (w WebhookMCTransactionSettledPayload) GetCounterpartAccountIdentifier() string {
	if len(w.UltimateCreditorAccountIdentifiers) == 0 {
		return ""
	}
	if w.InstructedAmount > 0 {
		return w.UltimateDebtorAccountIdentifiers[0].Identifier
	}
	return w.UltimateCreditorAccountIdentifiers[0].Identifier
}

func (w WebhookMCTransactionSettledPayload) GetCounterpartAccountOwner() string {
	if w.InstructedAmount > 0 {
		return w.UltimateDebtorName
	}
	return w.UltimateCreditorName
}

var _ Transaction = (*WebhookMCTransactionCancelledPayload)(nil)

// WebhookMCTransactionCancelledPayload
// This webhook confirms that a multicurrency transaction has been canceled
type WebhookMCTransactionCancelledPayload struct {
	TransactionID                      uuid.UUID        `json:"TransactionId" validate:"required"`
	BatchID                            uuid.UUID        `json:"BatchId" validate:"required"`
	EndToEndID                         string           `json:"EndToEndId" validate:"required"`
	SchemeEndToEndID                   string           `json:"SchemeEndToEndId"`
	AccountID                          uuid.UUID        `json:"AccountId" validate:"required"`
	Reference                          string           `json:"Reference" validate:"required"`
	UltimateCreditorName               string           `json:"UltimateCreditorName" validate:"required"`
	UltimateCreditorAccountIdentifiers []MCCYIdentifier `json:"UltimateCreditorAccountIdentifiers" validate:"required"`
	UltimateDebtorName                 string           `json:"UltimateDebtorName" validate:"required"`
	UltimateDebtorAccountIdentifiers   []MCCYIdentifier `json:"UltimateDebtorAccountIdentifiers" validate:"required"`
	InstructedAmount                   float64          `json:"InstructedAmount" validate:"required"`
	InstructedCurrency                 string           `json:"InstructedCurrency" validate:"required"`
	SchemePaymentMethod                string           `json:"SchemePaymentMethod" validate:"required"`
	TimestampCreated                   clearbank.Time   `json:"TimestampCreated" validate:"required"`
	TimestampCancelled                 clearbank.Time   `json:"TimestampCancelled" validate:"required"`
	CancellationCode                   string           `json:"CancellationCode" validate:"required"`
	CancellationReason                 string           `json:"CancellationReason" validate:"required"`
}

func (w WebhookMCTransactionCancelledPayload) GetID() uuid.UUID {
	return w.TransactionID
}

func (w WebhookMCTransactionCancelledPayload) GetEndToEndID() string {
	return w.EndToEndID
}

func (w WebhookMCTransactionCancelledPayload) GetCurrency() string {
	return w.InstructedCurrency
}

func (w WebhookMCTransactionCancelledPayload) GetAmount() float64 {
	return w.InstructedAmount
}

func (w WebhookMCTransactionCancelledPayload) IsReturn() bool {
	return false
}

func (w WebhookMCTransactionCancelledPayload) GetReference() string {
	return w.Reference
}

func (w WebhookMCTransactionCancelledPayload) GetAccountIdentifier() string {
	if len(w.UltimateCreditorAccountIdentifiers) == 0 {
		return ""
	}
	if w.InstructedAmount > 0 {
		return w.UltimateCreditorAccountIdentifiers[0].Identifier
	}
	return w.UltimateDebtorAccountIdentifiers[0].Identifier
}

func (w WebhookMCTransactionCancelledPayload) GetAccountOwner() string {
	if w.InstructedAmount > 0 {
		return w.UltimateCreditorName
	}
	return w.UltimateDebtorName
}

func (w WebhookMCTransactionCancelledPayload) GetCounterpartAccountIdentifier() string {
	if len(w.UltimateCreditorAccountIdentifiers) == 0 {
		return ""
	}
	if w.InstructedAmount > 0 {
		return w.UltimateDebtorAccountIdentifiers[0].Identifier
	}
	return w.UltimateCreditorAccountIdentifiers[0].Identifier
}

func (w WebhookMCTransactionCancelledPayload) GetCounterpartAccountOwner() string {
	if w.InstructedAmount > 0 {
		return w.UltimateDebtorName
	}
	return w.UltimateCreditorName
}

// WebhookMCPayloadAssessmentFailedPayload
// This webhook confirms that a multicurrency payment has failed assessment
type WebhookMCPayloadAssessmentFailedPayload struct {
	TransactionID                   uuid.UUID      `json:"TransactionId" validate:"required"`
	BatchID                         uuid.UUID      `json:"BatchId" validate:"required"`
	EndToEndID                      string         `json:"EndToEndId" validate:"required"`
	SchemeEndToEndID                string         `json:"SchemeEndToEndId"`
	InstructedAmount                float64        `json:"InstructedAmount" validate:"required"`
	InstructedCurrency              string         `json:"InstructedCurrency" validate:"required"`
	Error                           string         `json:"Error" validate:"required"`
	UltimateDebtorName              string         `json:"UltimateDebtorName"`
	UltimateDebtorAccount           string         `json:"UltimateDebtorAccount" validate:"required"`
	UltimateDebtorAccountIdentifier MCCYIdentifier `json:"UltimateDebtorAccountIdentifiers" validate:"required"`
	UltimateDebtorAddressLine1      string         `json:"UltimateDebtorAddressLine1" validate:"required"`
	UltimateDebtorAddressLine2      string         `json:"UltimateDebtorAddressLine2"`
	UltimateDebtorAddressLine3      string         `json:"UltimateDebtorAddressLine3"`
	UltimateDebtorPostCode          string         `json:"UltimateDebtorPostCode" validate:"required"`
	UltimateDebtorCountryCode       string         `json:"UltimateDebtorCountryCode" validate:"required"`
	UltimateCreditorIBAN            string         `json:"UltimateCreditorIBAN"`
	UltimateCreditorAccountNumber   string         `json:"UltimateCreditorAccountNumber"`
	UltimateCreditorBic             string         `json:"UltimateCreditorBic"`
	UltimateCreditorName            string         `json:"UltimateCreditorName" validate:"required"`
	UltimateCreditorAddressLine1    string         `json:"UltimateCreditorAddressLine1" validate:"required"`
	UltimateCreditorAddressLine2    string         `json:"UltimateCreditorAddressLine2" `
	UltimateCreditorAddressLine3    string         `json:"UltimateCreditorAddressLine3"`
	UltimateCreditorPostCode        string         `json:"UltimateCreditorPostCode" validate:"required"`
	UltimateCreditorCountryCode     string         `json:"UltimateCreditorCountryCode" validate:"required"`
}

// WebhookMCPayloadValidationFailedPayload
// This webhook confirms that a multicurrency payment has failed validation
type WebhookMCPayloadValidationFailedPayload struct {
	TransactionID                   uuid.UUID      `json:"TransactionId" validate:"required"`
	BatchID                         uuid.UUID      `json:"BatchId" validate:"required"`
	EndToEndID                      string         `json:"EndToEndId" validate:"required"`
	SchemeEndToEndID                string         `json:"SchemeEndToEndId"`
	InstructedAmount                float64        `json:"InstructedAmount" validate:"required"`
	InstructedCurrency              string         `json:"InstructedCurrency" validate:"required"`
	Errors                          []string       `json:"Errors" validate:"required"`
	UltimateDebtorName              string         `json:"UltimateDebtorName"`
	UltimateDebtorAccount           string         `json:"UltimateDebtorAccount" validate:"required"`
	UltimateDebtorAccountIdentifier MCCYIdentifier `json:"UltimateDebtorAccountIdentifiers" validate:"required"`
	UltimateDebtorAddressLine1      string         `json:"UltimateDebtorAddressLine1" validate:"required"`
	UltimateDebtorAddressLine2      string         `json:"UltimateDebtorAddressLine2"`
	UltimateDebtorAddressLine3      string         `json:"UltimateDebtorAddressLine3"`
	UltimateDebtorPostCode          string         `json:"UltimateDebtorPostCode" validate:"required"`
	UltimateDebtorCountryCode       string         `json:"UltimateDebtorCountryCode" validate:"required"`
	UltimateCreditorIBAN            string         `json:"UltimateCreditorIBAN"`
	UltimateCreditorAccountNumber   string         `json:"UltimateCreditorAccountNumber"`
	UltimateCreditorBic             string         `json:"UltimateCreditorBic"`
	UltimateCreditorName            string         `json:"UltimateCreditorName" validate:"required"`
	UltimateCreditorAddressLine1    string         `json:"UltimateCreditorAddressLine1" validate:"required"`
	UltimateCreditorAddressLine2    string         `json:"UltimateCreditorAddressLine2" `
	UltimateCreditorAddressLine3    string         `json:"UltimateCreditorAddressLine3"`
	UltimateCreditorPostCode        string         `json:"UltimateCreditorPostCode" validate:"required"`
	UltimateCreditorCountryCode     string         `json:"UltimateCreditorCountryCode" validate:"required"`
}
