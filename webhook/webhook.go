package webhook

import (
	"encoding/json"

	"github.com/google/uuid"
)

const (
	Test                           = "FITestEvent"
	AccountCreated                 = "AccountCreated"
	AccountDisabled                = "AccountDisabled"
	VirtualAccountCreated          = "VirtualAccountCreated"
	VirtualAccountCreationFailed   = "VirtualAccountCreationFailed"
	TransactionSettled             = "TransactionSettled"
	PaymentMessageAssesmentFailed  = "PaymentMessageAssesmentFailed"
	PaymentMessageValidationFailed = "PaymentMessageValidationFailed"
	TransactionRejected            = "TransactionRejected"
	MCCYAccountCreated             = "Accounts.AccountCreated"
	MCCYAccountUpdated             = "Accounts.AccountUpdated"
	MCCYVirtualAccountCreated      = "Accounts.VirtualAccountCreated"
	MCCYVirtualAccountUpdated      = "Accounts.VirtualAccountUpdated"
	MCCYInstitutionStatement       = "MCCY.Statements.InstitutionStatement"
	MCCYAccountStatement           = "MCCY.Statements.AccountStatement"
	MCCYTransactionCreated         = "Payments.Mccy.TransactionCreated"
	MCCYTransactionSettled         = "Payments.Mccy.TransactionSettled"
	MCCYTransactionCancelled       = "Payments.Mccy.TransactionCancelled"
	MCCYPayloadAssessmentFailed    = "Payments.Mccy.PaymentAssessmentFailed"
	MCCYPayloadValidationFailed    = "Payments.Mccy.PaymentValidationFailed"
	MCCYInternalTransfersSettled   = "Mccy.InternalTransfers.Settled"
	MCCYInternalTransfersCancelled = "Mccy.InternalTransfers.Cancelled"
	FxTradeExecuted                = "Fx.Trade.Executed"
	FxTradeSettled                 = "Fx.Trade.Settled"
	FxTradeCancelled               = "Fx.Trade.Cancelled"
)

type Request struct {
	Type    string          `json:"Type"`
	Version int             `json:"Version"`
	Nonce   int             `json:"Nonce"`
	Payload json.RawMessage `json:"Payload"`
}

type Response struct {
	// The value that you receive in the webhook request.
	Nonce int `json:"Nonce"`
}

// Transaction is a common interface for the transaction related webhooks.
type Transaction interface {
	GetID() uuid.UUID
	GetEndToEndID() string
	GetCurrency() string
	GetAmount() float64
	IsReturn() bool
	GetReference() string
	GetAccountIdentifier() string
	GetAccountOwner() string
	GetCounterpartAccountIdentifier() string
	GetCounterpartAccountOwner() string
}
