package webhook

import (
	"github.com/brokeyourbike/clearbank-api-client-go"
	"github.com/google/uuid"
)

type AdditionalProperty struct {
	Key   string `json:"Key" validate:"required"`
	Value string `json:"Value" validate:"required"`
}

// MCCYTransactionCreatedPayload
// This webhook confirms that a multicurrency transaction has been created
type MCCYTransactionCreatedPayload struct {
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

var _ Transaction = (*MCCYTransactionSettledPayload)(nil)

// MCCYTransactionSettledPayload
// This webhook confirms that a multicurrency transaction has been settled
type MCCYTransactionSettledPayload struct {
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

func (w MCCYTransactionSettledPayload) GetID() uuid.UUID {
	return w.TransactionID
}

func (w MCCYTransactionSettledPayload) GetEndToEndID() string {
	return w.EndToEndID
}

func (w MCCYTransactionSettledPayload) GetCurrency() string {
	return w.InstructedCurrency
}

func (w MCCYTransactionSettledPayload) GetAmount() float64 {
	return w.InstructedAmount
}

func (w MCCYTransactionSettledPayload) IsReturn() bool {
	return false
}

func (w MCCYTransactionSettledPayload) GetReference() string {
	return w.Reference
}

func (w MCCYTransactionSettledPayload) GetAccountIdentifier() string {
	if len(w.UltimateCreditorAccountIdentifiers) == 0 {
		return ""
	}
	if w.InstructedAmount > 0 {
		return w.UltimateCreditorAccountIdentifiers[0].Identifier
	}
	return w.UltimateDebtorAccountIdentifiers[0].Identifier
}

func (w MCCYTransactionSettledPayload) GetAccountOwner() string {
	if w.InstructedAmount > 0 {
		return w.UltimateCreditorName
	}
	return w.UltimateDebtorName
}

func (w MCCYTransactionSettledPayload) GetCounterpartAccountIdentifier() string {
	if len(w.UltimateCreditorAccountIdentifiers) == 0 {
		return ""
	}
	if w.InstructedAmount > 0 {
		return w.UltimateDebtorAccountIdentifiers[0].Identifier
	}
	return w.UltimateCreditorAccountIdentifiers[0].Identifier
}

func (w MCCYTransactionSettledPayload) GetCounterpartAccountOwner() string {
	if w.InstructedAmount > 0 {
		return w.UltimateDebtorName
	}
	return w.UltimateCreditorName
}

var _ Transaction = (*MCCYTransactionCancelledPayload)(nil)

// MCCYTransactionCancelledPayload
// This webhook confirms that a multicurrency transaction has been canceled
type MCCYTransactionCancelledPayload struct {
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

func (w MCCYTransactionCancelledPayload) GetID() uuid.UUID {
	return w.TransactionID
}

func (w MCCYTransactionCancelledPayload) GetEndToEndID() string {
	return w.EndToEndID
}

func (w MCCYTransactionCancelledPayload) GetCurrency() string {
	return w.InstructedCurrency
}

func (w MCCYTransactionCancelledPayload) GetAmount() float64 {
	return w.InstructedAmount
}

func (w MCCYTransactionCancelledPayload) IsReturn() bool {
	return false
}

func (w MCCYTransactionCancelledPayload) GetReference() string {
	return w.Reference
}

func (w MCCYTransactionCancelledPayload) GetAccountIdentifier() string {
	if len(w.UltimateCreditorAccountIdentifiers) == 0 {
		return ""
	}
	if w.InstructedAmount > 0 {
		return w.UltimateCreditorAccountIdentifiers[0].Identifier
	}
	return w.UltimateDebtorAccountIdentifiers[0].Identifier
}

func (w MCCYTransactionCancelledPayload) GetAccountOwner() string {
	if w.InstructedAmount > 0 {
		return w.UltimateCreditorName
	}
	return w.UltimateDebtorName
}

func (w MCCYTransactionCancelledPayload) GetCounterpartAccountIdentifier() string {
	if len(w.UltimateCreditorAccountIdentifiers) == 0 {
		return ""
	}
	if w.InstructedAmount > 0 {
		return w.UltimateDebtorAccountIdentifiers[0].Identifier
	}
	return w.UltimateCreditorAccountIdentifiers[0].Identifier
}

func (w MCCYTransactionCancelledPayload) GetCounterpartAccountOwner() string {
	if w.InstructedAmount > 0 {
		return w.UltimateDebtorName
	}
	return w.UltimateCreditorName
}

// MCCYPayloadAssessmentFailedPayload
// This webhook confirms that a multicurrency payment has failed assessment
type MCCYPayloadAssessmentFailedPayload struct {
	TransactionID                   uuid.UUID `json:"TransactionId" validate:"required"`
	BatchID                         uuid.UUID `json:"BatchId" validate:"required"`
	EndToEndID                      string    `json:"EndToEndId" validate:"required"`
	SchemeEndToEndID                string    `json:"SchemeEndToEndId"`
	InstructedAmount                float64   `json:"InstructedAmount" validate:"required"`
	InstructedCurrency              string    `json:"InstructedCurrency" validate:"required"`
	Error                           string    `json:"Error" validate:"required"`
	UltimateDebtorName              string    `json:"UltimateDebtorName"`
	UltimateDebtorAccount           string    `json:"UltimateDebtorAccount" validate:"required"`
	UltimateDebtorAccountIdentifier string    `json:"UltimateDebtorAccountIdentifier" validate:"required"`
	UltimateDebtorAddressLine1      string    `json:"UltimateDebtorAddressLine1" validate:"required"`
	UltimateDebtorAddressLine2      string    `json:"UltimateDebtorAddressLine2"`
	UltimateDebtorAddressLine3      string    `json:"UltimateDebtorAddressLine3"`
	UltimateDebtorPostCode          string    `json:"UltimateDebtorPostCode" validate:"required"`
	UltimateDebtorCountryCode       string    `json:"UltimateDebtorCountryCode" validate:"required"`
	UltimateCreditorIBAN            string    `json:"UltimateCreditorIBAN"`
	UltimateCreditorAccountNumber   string    `json:"UltimateCreditorAccountNumber"`
	UltimateCreditorBic             string    `json:"UltimateCreditorBic"`
	UltimateCreditorName            string    `json:"UltimateCreditorName" validate:"required"`
	UltimateCreditorAddressLine1    string    `json:"UltimateCreditorAddressLine1" validate:"required"`
	UltimateCreditorAddressLine2    string    `json:"UltimateCreditorAddressLine2" `
	UltimateCreditorAddressLine3    string    `json:"UltimateCreditorAddressLine3"`
	UltimateCreditorPostCode        string    `json:"UltimateCreditorPostCode" validate:"required"`
	UltimateCreditorCountryCode     string    `json:"UltimateCreditorCountryCode" validate:"required"`
}

// MCCYPayloadValidationFailedPayload
// This webhook confirms that a multicurrency payment has failed validation
type MCCYPayloadValidationFailedPayload struct {
	TransactionID                   uuid.UUID `json:"TransactionId" validate:"required"`
	BatchID                         uuid.UUID `json:"BatchId" validate:"required"`
	EndToEndID                      string    `json:"EndToEndId" validate:"required"`
	SchemeEndToEndID                string    `json:"SchemeEndToEndId"`
	InstructedAmount                float64   `json:"InstructedAmount" validate:"required"`
	InstructedCurrency              string    `json:"InstructedCurrency" validate:"required"`
	Errors                          []string  `json:"Errors" validate:"required"`
	UltimateDebtorName              string    `json:"UltimateDebtorName"`
	UltimateDebtorAccount           string    `json:"UltimateDebtorAccount" validate:"required"`
	UltimateDebtorAccountIdentifier string    `json:"UltimateDebtorAccountIdentifier" validate:"required"`
	UltimateDebtorAddressLine1      string    `json:"UltimateDebtorAddressLine1" validate:"required"`
	UltimateDebtorAddressLine2      string    `json:"UltimateDebtorAddressLine2"`
	UltimateDebtorAddressLine3      string    `json:"UltimateDebtorAddressLine3"`
	UltimateDebtorPostCode          string    `json:"UltimateDebtorPostCode" validate:"required"`
	UltimateDebtorCountryCode       string    `json:"UltimateDebtorCountryCode" validate:"required"`
	UltimateCreditorIBAN            string    `json:"UltimateCreditorIBAN"`
	UltimateCreditorAccountNumber   string    `json:"UltimateCreditorAccountNumber"`
	UltimateCreditorBic             string    `json:"UltimateCreditorBic"`
	UltimateCreditorName            string    `json:"UltimateCreditorName" validate:"required"`
	UltimateCreditorAddressLine1    string    `json:"UltimateCreditorAddressLine1" validate:"required"`
	UltimateCreditorAddressLine2    string    `json:"UltimateCreditorAddressLine2" `
	UltimateCreditorAddressLine3    string    `json:"UltimateCreditorAddressLine3"`
	UltimateCreditorPostCode        string    `json:"UltimateCreditorPostCode" validate:"required"`
	UltimateCreditorCountryCode     string    `json:"UltimateCreditorCountryCode" validate:"required"`
}
