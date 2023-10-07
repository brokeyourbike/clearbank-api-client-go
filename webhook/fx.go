package webhook

import (
	"github.com/brokeyourbike/clearbank-api-client-go"
	"github.com/google/uuid"
)

// FxTradeExecutedPayload
// This webhook confirms that the FX trade has been executed
type FxTradeExecutedPayload struct {
	SellAccountOwner        string              `json:"SellAccountOwner" validate:"required"`
	SellAccountIBAN         string              `json:"SellAccountIban" validate:"required"`
	BuyAccountOwner         string              `json:"BuyAccountOwner" validate:"required"`
	BuyAccountIBAN          string              `json:"BuyAccountIban" validate:"required"`
	Attestation             string              `json:"Attestation" validate:"required"`
	ValueDate               clearbank.Time      `json:"ValueDate" validate:"required"`
	InstructedAmount        float64             `json:"InstructedAmount" validate:"required"`
	FixedSide               clearbank.FixedSide `json:"FixedSide" validate:"required"`
	BuyCurrency             string              `json:"BuyCurrency" validate:"required"`
	SellCurrency            string              `json:"SellCurrency" validate:"required"`
	Margin                  float64             `json:"Margin"`
	MarginAccountIBAN       string              `json:"MarginAccountIban"`
	EndToEndID              string              `json:"EndToEndId" validate:"required"`
	UnstructuredInformation string              `json:"UnstructuredInformation"`
	Symbol                  string              `json:"Symbol" validate:"required"`
	FIReceivedMarginAmount  float64             `json:"FIReceivedMarginAmount"`
	BuyAmount               float64             `json:"BuyAmount" validate:"required"`
	SellAmount              float64             `json:"SellAmount" validate:"required"`
	ClearbankRate           float64             `json:"ClearBankRate" validate:"required"`
	FIRate                  float64             `json:"FIRate" validate:"required"`
	RequestedTime           clearbank.Time      `json:"RequestedTime" validate:"required"`
	ExecutedTime            clearbank.Time      `json:"ExecutedTime" validate:"required"`
}

// FxTradeSettledPayload
// This webhook confirms that the FX trade has been settled
type FxTradeSettledPayload struct {
	SellAccountOwner        string              `json:"SellAccountOwner" validate:"required"`
	SellAccountIBAN         string              `json:"SellAccountIban" validate:"required"`
	BuyAccountOwner         string              `json:"BuyAccountOwner" validate:"required"`
	BuyAccountIBAN          string              `json:"BuyAccountIban" validate:"required"`
	Attestation             string              `json:"Attestation" validate:"required"`
	ValueDate               clearbank.Time      `json:"ValueDate" validate:"required"`
	InstructedAmount        float64             `json:"InstructedAmount" validate:"required"`
	FixedSide               clearbank.FixedSide `json:"FixedSide" validate:"required"`
	BuyCurrency             string              `json:"BuyCurrency" validate:"required"`
	SellCurrency            string              `json:"SellCurrency" validate:"required"`
	Margin                  float64             `json:"Margin"`
	MarginAccountIBAN       string              `json:"MarginAccountIban"`
	EndToEndID              string              `json:"EndToEndId" validate:"required"`
	UnstructuredInformation string              `json:"UnstructuredInformation"`
	DebitTransactionID      uuid.UUID           `json:"DebitTransactionId" validate:"required"`
	DebitAmount             float64             `json:"DebitAmount" validate:"required"`
	DebitCurrency           string              `json:"DebitCurrency" validate:"required"`
	DebitCreatedAt          clearbank.Time      `json:"DebitCreatedAt" validate:"required"`
	DebitSettledAt          clearbank.Time      `json:"DebitSettledAt" validate:"required"`
	CreditTransactionID     uuid.UUID           `json:"CreditTransactionId" validate:"required"`
	CreditAmount            float64             `json:"CreditAmount" validate:"required"`
	CreditCurrency          string              `json:"CreditCurrency" validate:"required"`
	CreditCreatedAt         clearbank.Time      `json:"CreditCreatedAt" validate:"required"`
	CreditSettledAt         clearbank.Time      `json:"CreditSettledAt" validate:"required"`
	MarginTransactionID     *uuid.UUID          `json:"MarginTransactionId"`
	MarginAmount            float64             `json:"MarginAmount"`
	MarginCurrency          string              `json:"MarginCurrency"`
	MarginCreatedAt         string              `json:"MarginCreatedAt"`
	MarginSettledAt         string              `json:"MarginSettledAt"`
}

// FxTradeCancelledPayload
// This webhook confirms that the FX trade has been canceled
type FxTradeCancelledPayload struct {
	SellAccountOwner        string         `json:"SellAccountOwner" validate:"required"`
	SellAccountIBAN         string         `json:"SellAccountIban" validate:"required"`
	BuyAccountOwner         string         `json:"BuyAccountOwner" validate:"required"`
	BuyAccountIBAN          string         `json:"BuyAccountIban" validate:"required"`
	Attestation             string         `json:"Attestation" validate:"required"`
	ValueDate               clearbank.Time `json:"ValueDate" validate:"required"`
	InstructedAmount        float64        `json:"InstructedAmount" validate:"required"`
	FixedSide               string         `json:"FixedSide" validate:"required"`
	BuyCurrency             string         `json:"BuyCurrency" validate:"required"`
	SellCurrency            string         `json:"SellCurrency" validate:"required"`
	Margin                  float64        `json:"Margin"`
	MarginAccountIBAN       string         `json:"MarginAccountIban"`
	EndToEndID              string         `json:"EndToEndId" validate:"required"`
	UnstructuredInformation string         `json:"UnstructuredInformation"`
	CancellationCode        string         `json:"CancellationCode" validate:"required"`
	CancellationReason      string         `json:"CancellationReason" validate:"required"`
}
