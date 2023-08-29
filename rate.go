package clearbank

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type RateClient interface {
	FetchMarketrate(context.Context, MarketrateParams) (MarketrateResponse, error)
	Negotiate(context.Context) (NegotiateResponse, error)
}

type FixedSide string

const (
	FixedSideBuy  FixedSide = "BUY"
	FixedSideSell FixedSide = "SELL"
)

type MarketrateParams struct {
	SellCurrency     string
	BuyCurrency      string
	InstructedAmount float64
	FixedSide        FixedSide
	ValueDate        string
}

type MarketrateResponse struct {
	Symbol     string  `json:"symbol"`
	MarketRate float64 `json:"marketRate"`
	ValueDate  Time    `json:"valueDate"`
}

func (c *client) FetchMarketrate(ctx context.Context, params MarketrateParams) (data MarketrateResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/fx/v1/marketrate", nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.AddQueryParams(map[string]string{
		"SellCurrency":     params.SellCurrency,
		"BuyCurrency":      params.BuyCurrency,
		"InstructedAmount": strconv.FormatFloat(params.InstructedAmount, 'f', 2, 64),
		"FixedSide":        string(params.FixedSide),
		"ValueDate":        params.ValueDate,
	})

	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)

	return data, c.do(ctx, req)
}

type NegotiateResponse struct {
	URL string `json:"url"`
}

func (c *client) Negotiate(ctx context.Context) (data NegotiateResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/fx/v1/marketdatastream/negotiate", nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)

	return data, c.do(ctx, req)
}
