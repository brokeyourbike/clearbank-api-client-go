package webhook

import (
	"github.com/brokeyourbike/clearbank-api-client-go"
	"github.com/google/uuid"
)

type TransactionStatus string

const (
	TransactionStatusSettled  TransactionStatus = "Settled"
	TransactionStatusRejected TransactionStatus = "Rejected"
)

type TransactionScheme string

const (
	TransactionSchemeTransfer       = "Transfer"
	TransactionSchemeFasterPayments = "FasterPayments"
	TransactionSchemeBacs           = "Bacs"
	TransactionSchemeChaps          = "Chaps"
)

type DebitCreditCode string

const (
	DebitCreditCodeDebit  DebitCreditCode = "Debit"
	DebitCreditCodeCredit DebitCreditCode = "Credit"
)

type TransactionAccount struct {
	IBAN                 string `json:"IBAN"`
	BBAN                 string `json:"BBAN"`
	OwnerName            string `json:"OwnerName"`
	TransactionOwnerName string `json:"TransactionOwnerName"`
	InstitutionName      string `json:"InstitutionName"`
}

var _ Transaction = (*TransactionSettledPayload)(nil)

// TransactionSettledPayload
// This sends a webhook notification confirming the transaction has settled
type TransactionSettledPayload struct {
	TransactionID               uuid.UUID          `json:"TransactionId" validate:"required"`
	Status                      TransactionStatus  `json:"Status" validate:"required"`
	Scheme                      TransactionScheme  `json:"Scheme" validate:"required"`
	EndToEndTransactionID       string             `json:"EndToEndTransactionId" validate:"required"`
	Amount                      float64            `json:"Amount" validate:"required"`
	TimestampSettled            clearbank.Time     `json:"TimestampSettled" validate:"required"`
	TimestampCreated            clearbank.Time     `json:"TimestampCreated" validate:"required"`
	CurrencyCode                string             `json:"CurrencyCode" validate:"required"`
	DebitCreditCode             DebitCreditCode    `json:"DebitCreditCode" validate:"required"`
	Reference                   string             `json:"Reference"`
	Return                      bool               `json:"IsReturn" validate:"boolean"`
	Account                     TransactionAccount `json:"Account" validate:"required"`
	CounterpartAccount          TransactionAccount `json:"CounterpartAccount" validate:"required"`
	ActualEndToEndTransactionID string             `json:"ActualEndToEndTransactionId" validate:"required"`
	DirectDebitMandateID        string             `json:"DirectDebitMandateId"`
	TransactionCode             string             `json:"TransactionCode"`
	ServiceUserNumber           string             `json:"ServiceUserNumber"`
	BacsTransactionID           string             `json:"BacsTransactionId"`
	BacsTransactionDescription  string             `json:"BacsTransactionDescription"`
	TransactionType             string             `json:"TransactionType"`
	TransactionSource           string             `json:"TransactionSource"`
	SupplementaryData           []struct {
		Name  string `json:"Name"`
		Value string `json:"Value"`
	} `json:"SupplementaryData"`
}

func (w TransactionSettledPayload) GetID() uuid.UUID {
	return w.TransactionID
}

func (w TransactionSettledPayload) GetEndToEndID() string {
	return w.EndToEndTransactionID
}

func (w TransactionSettledPayload) GetCurrency() string {
	return w.CurrencyCode
}

func (w TransactionSettledPayload) GetAmount() float64 {
	if w.DebitCreditCode == DebitCreditCodeDebit {
		return -w.Amount
	}
	return w.Amount
}

func (w TransactionSettledPayload) IsReturn() bool {
	return w.Return
}

func (w TransactionSettledPayload) GetReference() string {
	return w.Reference
}

func (w TransactionSettledPayload) GetAccountIdentifier() string {
	return w.Account.IBAN
}

func (w TransactionSettledPayload) GetAccountOwner() string {
	return w.Account.OwnerName
}

func (w TransactionSettledPayload) GetCounterpartAccountIdentifier() string {
	return w.CounterpartAccount.IBAN
}

func (w TransactionSettledPayload) GetCounterpartAccountOwner() string {
	return w.CounterpartAccount.TransactionOwnerName
}

var _ Transaction = (*TransactionRejectedPayload)(nil)

// TransactionRejectedPayload
// This webhook confirms the payment has been rejected
type TransactionRejectedPayload struct {
	TransactionID         uuid.UUID          `json:"TransactionId" validate:"required"`
	Status                TransactionStatus  `json:"Status" validate:"required"`
	Scheme                TransactionScheme  `json:"Scheme" validate:"required"`
	EndToEndTransactionID string             `json:"EndToEndTransactionId" validate:"required"`
	Amount                float64            `json:"Amount" validate:"required"`
	TimestampModified     clearbank.Time     `json:"TimestampModified" validate:"required"`
	CurrencyCode          string             `json:"CurrencyCode"`
	DebitCreditCode       DebitCreditCode    `json:"DebitCreditCode" validate:"required"`
	Reference             string             `json:"Reference"`
	Return                bool               `json:"IsReturn" validate:"boolean"`
	CancellationReason    string             `json:"CancellationReason"`
	CancellationCode      string             `json:"CancellationCode"`
	Account               TransactionAccount `json:"Account" validate:"required"`
	CounterpartAccount    TransactionAccount `json:"CounterpartAccount" validate:"required"`
}

func (w TransactionRejectedPayload) GetID() uuid.UUID {
	return w.TransactionID
}

func (w TransactionRejectedPayload) GetEndToEndID() string {
	return w.EndToEndTransactionID
}

func (w TransactionRejectedPayload) GetCurrency() string {
	return w.CurrencyCode
}

func (w TransactionRejectedPayload) GetAmount() float64 {
	if w.DebitCreditCode == DebitCreditCodeDebit {
		return -w.Amount
	}
	return w.Amount
}

func (w TransactionRejectedPayload) IsReturn() bool {
	return w.Return
}

func (w TransactionRejectedPayload) GetReference() string {
	return w.Reference
}

func (w TransactionRejectedPayload) GetAccountIdentifier() string {
	return w.Account.IBAN
}

func (w TransactionRejectedPayload) GetAccountOwner() string {
	return w.Account.OwnerName
}

func (w TransactionRejectedPayload) GetCounterpartAccountIdentifier() string {
	return w.CounterpartAccount.IBAN
}

func (w TransactionRejectedPayload) GetCounterpartAccountOwner() string {
	return w.CounterpartAccount.TransactionOwnerName
}

// PaymentMessageAssesmentFailedPayload
// This webhook confirms the payment assessment has failed
type PaymentMessageAssesmentFailedPayload struct {
	MessageID         uuid.UUID         `json:"MessageId" validate:"required"`
	PaymentMethodType TransactionScheme `json:"PaymentMethodType" validate:"required"`
	AssesmentFailures []struct {
		EndToEndID string   `json:"EndToEndId"`
		Reasons    []string `json:"Reasons"`
	} `json:"AssesmentFailure" validate:"required"`
	AccountIdentification struct {
		Debtor struct {
			IBAN string `json:"IBAN" validate:"required"`
			BBAN string `json:"BBAN" validate:"required"`
		}
		Creditors []struct {
			Reference    string  `json:"Reference"`
			Amount       float64 `json:"Amount"`
			CurrencyCode string  `json:"CurrencyCode"`
			IBAN         string  `json:"IBAN" validate:"required"`
			BBAN         string  `json:"BBAN" validate:"required"`
		} `json:"Creditors" validate:"required"`
	} `json:"AccountIdentification" validate:"required"`
}

// PaymentMessageValidationFailedPayload
// This webhook confirms the payment validation has failed
type PaymentMessageValidationFailedPayload struct {
	MessageID          uuid.UUID         `json:"MessageId" validate:"required"`
	PaymentMethodType  TransactionScheme `json:"PaymentMethodType" validate:"required"`
	ValidationFailures []struct {
		EndToEndID string   `json:"EndToEndId"`
		Reasons    []string `json:"Reasons"`
	} `json:"ValidationFailure" validate:"required"`
	AccountIdentification struct {
		Debtor struct {
			IBAN               string `json:"IBAN" validate:"required"`
			BBAN               string `json:"BBAN" validate:"required"`
			AccountName        string `json:"AccountName"`
			AccountHolderLabel string `json:"AccountHolderLabel"`
			InstitutionName    string `json:"InstitutionName"`
		}
		Creditors []struct {
			Reference          string  `json:"Reference" validate:"required"`
			Amount             float64 `json:"Amount" validate:"required"`
			CurrencyCode       string  `json:"CurrencyCode" validate:"required"`
			IBAN               string  `json:"IBAN" validate:"required"`
			BBAN               string  `json:"BBAN" validate:"required"`
			AccountName        string  `json:"AccountName"`
			AccountHolderLabel string  `json:"AccountHolderLabel"`
			InstitutionName    string  `json:"InstitutionName"`
		} `json:"Creditors" validate:"required"`
	} `json:"AccountIdentification" validate:"required"`
}

// InboundHeldTransactionPayload
// This webhook confirms the inbound transaction has been held
type InboundHeldTransactionPayload struct {
	Scheme  string `json:"Scheme" validate:"required"`
	Account struct {
		IBAN string `json:"IBAN" validate:"required"`
		BBAN string `json:"BBAN" validate:"required"`
	} `json:"Account" validate:"required"`
	CounterpartAccount struct {
		IBAN string `json:"IBAN" validate:"required"`
		BBAN string `json:"BBAN" validate:"required"`
	} `json:"CounterpartAccount" validate:"required"`
	TransactionAmount     float64        `json:"TransactionAmount"`
	PaymentReference      string         `json:"PaymentReference"`
	EndToEndTransactionID string         `json:"EndToEndTransactionID"`
	TimestampCreated      clearbank.Time `json:"TimestampCreated"`
}

// OutboundHeldTransactionPayload
// This webhook confirms the outbound transaction has been held
type OutboundHeldTransactionPayload struct {
	Scheme  string `json:"Scheme" validate:"required"`
	Account struct {
		IBAN string `json:"IBAN" validate:"required"`
		BBAN string `json:"BBAN" validate:"required"`
	} `json:"Account" validate:"required"`
	CounterpartAccount struct {
		IBAN string `json:"IBAN" validate:"required"`
		BBAN string `json:"BBAN" validate:"required"`
	} `json:"CounterpartAccount" validate:"required"`
	TransactionAmount     float64        `json:"TransactionAmount"`
	PaymentReference      string         `json:"PaymentReference"`
	EndToEndTransactionID string         `json:"EndToEndTransactionID"`
	TimestampCreated      clearbank.Time `json:"TimestampCreated"`
}
