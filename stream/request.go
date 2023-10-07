package stream

import (
	"encoding/json"

	"github.com/brokeyourbike/clearbank-api-client-go"
)

type SubscriptionRequest struct {
	Type    string `json:"Type"`
	Version int    `json:"Version"`
	Payload struct {
		CurrencyPair string `json:"CurrencyPair"`
	} `json:"Payload"`
}

func NewSubscriptionRequest(symbol string) SubscriptionRequest {
	req := SubscriptionRequest{Type: "SubscriptionRequest", Version: 1}
	req.Payload.CurrencyPair = symbol
	return req
}

type MarketMessage struct {
	Type    string          `json:"Type"`
	Version int             `json:"Version"`
	Payload json.RawMessage `json:"Payload"`
}

type HeartbeatPayload struct {
	Time clearbank.Time `json:"Time"`
}

type MarketDataPayload struct {
	Symbol      string         `json:"Symbol"`
	Currency    string         `json:"Currency"`
	ValueDate   clearbank.Time `json:"ValueDate"`
	Type        string         `json:"Type"`
	SendingTime string         `json:"SendingTime"`
	Entries     []struct {
		Size  float64 `json:"Size"`
		Price float64 `json:"Price"`
	} `json:"Entries"`
}
