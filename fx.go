package clearbank

import (
	"context"
	"fmt"
	"net/http"
)

type FxClient interface {
	InitiateFxOrder(context.Context, FXPayload) error
}

type FXAttestation string

const (
	FXAttestationSameOwner      FXAttestation = "Y"
	FXAttestationDifferentOwner FXAttestation = "N"
)

type FXPayload struct {
	CustomerInformation struct {
		SellAccount struct {
			Owner string `json:"owner"`
			IBAN  string `json:"iban"`
		} `json:"sellAccount"`
		BuyAccount struct {
			Owner string `json:"owner"`
			IBAN  string `json:"iban"`
		} `json:"buyAccount"`
		Attestation FXAttestation `json:"attestation"`
	} `json:"customerInformation"`
	TradeInformation struct {
		ValueDate string `json:"valueDate,omitempty"`
		Details   struct {
			InstructedAmount float64   `json:"instructedAmount"`
			FixedSide        FixedSide `json:"fixedSide"`
			SellCurrency     string    `json:"sellCurrency"`
			BuyCurrency      string    `json:"buyCurrency"`
		} `json:"details"`
		Margin struct {
			Amount  float64 `json:"amount"`
			Account struct {
				Owner string `json:"owner"`
				IBAN  string `json:"iban"`
			} `json:"account"`
		} `json:"margin,omitempty"`
		EndToEndID              string `json:"endToEndId"`
		UnstructuredInformation string `json:"unstructuredInformation,omitempty"`
	} `json:"tradeInformation"`
}

func (c *client) InitiateFxOrder(ctx context.Context, payload FXPayload) error {
	req, err := c.newRequest(ctx, http.MethodPost, "/fx/v1/order", payload)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusAccepted)
	return c.do(ctx, req)
}
