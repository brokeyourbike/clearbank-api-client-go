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

var _ Transaction = (*WebhookTransactionSettledPayload)(nil)

// WebhookTransactionSettledPayload
// This sends a webhook notification confirming the transaction has settled
type WebhookTransactionSettledPayload struct {
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

func (w WebhookTransactionSettledPayload) GetID() uuid.UUID {
	return w.TransactionID
}

func (w WebhookTransactionSettledPayload) GetEndToEndID() string {
	return w.EndToEndTransactionID
}

func (w WebhookTransactionSettledPayload) GetCurrency() string {
	return w.CurrencyCode
}

func (w WebhookTransactionSettledPayload) GetAmount() float64 {
	if w.DebitCreditCode == DebitCreditCodeDebit {
		return -w.Amount
	}
	return w.Amount
}

func (w WebhookTransactionSettledPayload) IsReturn() bool {
	return w.Return
}

func (w WebhookTransactionSettledPayload) GetReference() string {
	return w.Reference
}

func (w WebhookTransactionSettledPayload) GetAccountIdentifier() string {
	return w.Account.IBAN
}

func (w WebhookTransactionSettledPayload) GetAccountOwner() string {
	return w.Account.OwnerName
}

func (w WebhookTransactionSettledPayload) GetCounterpartAccountIdentifier() string {
	return w.CounterpartAccount.IBAN
}

func (w WebhookTransactionSettledPayload) GetCounterpartAccountOwner() string {
	return w.CounterpartAccount.TransactionOwnerName
}

var _ Transaction = (*WebhookTransactionRejectedPayload)(nil)

// WebhookTransactionRejectedPayload
// This webhook confirms the payment has been rejected
type WebhookTransactionRejectedPayload struct {
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

func (w WebhookTransactionRejectedPayload) GetID() uuid.UUID {
	return w.TransactionID
}

func (w WebhookTransactionRejectedPayload) GetEndToEndID() string {
	return w.EndToEndTransactionID
}

func (w WebhookTransactionRejectedPayload) GetCurrency() string {
	return w.CurrencyCode
}

func (w WebhookTransactionRejectedPayload) GetAmount() float64 {
	if w.DebitCreditCode == DebitCreditCodeDebit {
		return -w.Amount
	}
	return w.Amount
}

func (w WebhookTransactionRejectedPayload) IsReturn() bool {
	return w.Return
}

func (w WebhookTransactionRejectedPayload) GetReference() string {
	return w.Reference
}

func (w WebhookTransactionRejectedPayload) GetAccountIdentifier() string {
	return w.Account.IBAN
}

func (w WebhookTransactionRejectedPayload) GetAccountOwner() string {
	return w.Account.OwnerName
}

func (w WebhookTransactionRejectedPayload) GetCounterpartAccountIdentifier() string {
	return w.CounterpartAccount.IBAN
}

func (w WebhookTransactionRejectedPayload) GetCounterpartAccountOwner() string {
	return w.CounterpartAccount.TransactionOwnerName
}

// WebhookPaymentMessageAssesmentFailedPayload
// This webhook confirms the payment assessment has failed
type WebhookPaymentMessageAssesmentFailedPayload struct {
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

// WebhookPaymentMessageValidationFailedPayload
// This webhook confirms the payment validation has failed
type WebhookPaymentMessageValidationFailedPayload struct {
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
