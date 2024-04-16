package clearbank

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type FxClient interface {
	InitiateFxOrder(context.Context, FXPayload) error
	RequestFxQuote(ctx context.Context, payload FXQuotePayload) (data FXQuoteResponse, err error)
	ExecuteFxQuote(ctx context.Context, payload FXPayload) error
}

type FXAttestation string

const (
	FXAttestationSameOwner      FXAttestation = "Y"
	FXAttestationDifferentOwner FXAttestation = "N"
)

type FXAccount struct {
	Owner string `json:"owner"`
	IBAN  string `json:"iban"`
}

type FXTradeInformation struct {
	ValueDate string `json:"valueDate,omitempty"`
	Details   struct {
		InstructedAmount float64   `json:"instructedAmount"`
		FixedSide        FixedSide `json:"fixedSide"`
		SellCurrency     string    `json:"sellCurrency"`
		BuyCurrency      string    `json:"buyCurrency"`
	} `json:"details"`
	Margin struct {
		Amount  float64   `json:"amount"`
		Account FXAccount `json:"account"`
	} `json:"margin,omitempty"`
	EndToEndID              string `json:"endToEndId"`
	UnstructuredInformation string `json:"unstructuredInformation,omitempty"`
}

type FXPayload struct {
	QuoteID             *uuid.UUID `json:"quoteId,omitempty"`
	CustomerInformation struct {
		SellAccount FXAccount     `json:"sellAccount"`
		BuyAccount  FXAccount     `json:"buyAccount"`
		Attestation FXAttestation `json:"attestation"`
	} `json:"customerInformation"`
	MarginAccount    *FXAccount          `json:"marginAccount,omitempty"`
	TradeInformation *FXTradeInformation `json:"tradeInformation,omitempty"`
}

func (c *client) InitiateFxOrder(ctx context.Context, payload FXPayload) error {
	req, err := c.newRequest(ctx, http.MethodPost, "/fx/v1/order", payload)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusAccepted)
	return c.do(ctx, req)
}

type QuoteFixedSide string

const (
	QuoteFixedSideBuy  QuoteFixedSide = "Buy"
	QuoteFixedSideSell QuoteFixedSide = "Sell"
)

type FXQuotePayload struct {
	SellCurrency     string         `json:"SellCurrency"`
	BuyCurrency      string         `json:"BuyCurrency"`
	InstructedAmount float64        `json:"InstructedAmount"`
	FixedSide        QuoteFixedSide `json:"FixedSide"`
	ValueDate        string         `json:"ValueDate"`
	Margin           *float64       `json:"Margin,omitempty"`
}

type FXQuoteResponse struct {
	QuoteID      uuid.UUID `json:"QuoteId"`
	ValueDate    string    `json:"ValueDate"`
	CurrencyPair string    `json:"CurrencyPair"`
	ExchangeRate float64   `json:"ExchangeRate"`
	SellCurrency string    `json:"SellCurrency"`
	SellAmount   float64   `json:"SellAmount"`
	BuyCurrency  string    `json:"BuyCurrency"`
	BuyAmount    float64   `json:"BuyAmount"`
	CreatedAt    Time      `json:"CreatedAt"`
	ExpiresAt    Time      `json:"ExpiresAt"`
	QuoteRequest struct {
		ValueDate        string  `json:"ValueDate"`
		InstructedAmount float64 `json:"InstructedAmount"`
		FixedSide        string  `json:"FixedSide"`
		SellCurrency     string  `json:"SellCurrency"`
		BuyCurrency      string  `json:"BuyCurrency"`
	} `json:"QuoteRequest"`
}

func (q FXQuoteResponse) GetRate() float64 {
	if strings.HasPrefix(q.CurrencyPair, q.SellCurrency) {
		return q.ExchangeRate
	}
	return 1 / q.ExchangeRate
}

func (c *client) RequestFxQuote(ctx context.Context, payload FXQuotePayload) (data FXQuoteResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/fx/v1/quote", payload)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)
	return data, c.do(ctx, req)
}

func (c *client) ExecuteFxQuote(ctx context.Context, payload FXPayload) (err error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/fx/v1/executeQuote", payload)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusAccepted)
	return c.do(ctx, req)
}
