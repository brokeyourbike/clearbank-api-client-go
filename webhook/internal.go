package webhook

// WebhookMCInternalTransfersSettledPayload
// This webhook confirms that a multicurrency internal transfer request has been settled
type WebhookMCInternalTransfersSettledPayload struct {
	EndToEndID         string  `json:"EndToEndId" validate:"required"`
	InstructedCurrency string  `json:"InstructedCurrency" validate:"required"`
	InstructedAmount   float64 `json:"InstructedAmount" validate:"required"`
	DebitAccountIBAN   string  `json:"DebitAccountIban"`
	CreditAccountIBAN  string  `json:"CreditAccountIban"`
	Reference          string  `json:"Reference"`
}

// WebhookMCInternalTransfersCancelledPayload
// This webhook confirms that a multicurrency internal transfer request has been canceled
type WebhookMCInternalTransfersCancelledPayload struct {
	CancellationCode   string   `json:"CancellationCode" validate:"required"`
	FailureReasons     []string `json:"FailureReasons" validate:"required"`
	EndToEndID         string   `json:"EndToEndId" validate:"required"`
	InstructedCurrency string   `json:"InstructedCurrency" validate:"required"`
	InstructedAmount   float64  `json:"InstructedAmount" validate:"required"`
	DebitAccountIBAN   string   `json:"DebitAccountIban"`
	CreditAccountIBAN  string   `json:"CreditAccountIban"`
	Reference          string   `json:"Reference"`
}
