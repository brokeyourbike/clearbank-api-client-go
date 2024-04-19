package clearbank

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type MCCYTransactionsClient interface {
	InitiateInternalTransaction(context.Context, CreateInternalTransactionPayload) error
	InitiateMCCYTransactions(context.Context, CreateMCCYTransactionsPayload) error
	InitiateMCCYInboundPayment(ctx context.Context, accountUniqueID string, payload CreateMCCYInboundPaymentPayload) (MCCYInboundPaymentResponse, error)
	FetchMCCYTransaction(ctx context.Context, trxID uuid.UUID) (MCCYTransactionResponse, error)
	FetchMCCYTransactionsForAccount(ctx context.Context, accountID uuid.UUID, currency string, params FetchTransactionsParams) (MCCYTransactionsResponse, error)
	FetchMCCYTransactionsForVirtualAccount(ctx context.Context, virtualAccountID uuid.UUID, currency string, params FetchTransactionsParams) (MCCYTransactionsResponse, error)
}

type CreateInternalTransactionPayload struct {
	DebitAccountIBAN   string  `json:"debitAccountIban"`
	CreditAccountIBAN  string  `json:"creditAccountIban"`
	InstructedAmount   float64 `json:"instructedAmount"`
	InstructedCurrency string  `json:"instructedCurrency"`
	EndToEndID         string  `json:"endToEndId"`
	Reference          string  `json:"reference"`
}

func (c *client) InitiateInternalTransaction(ctx context.Context, payload CreateInternalTransactionPayload) error {
	req, err := c.newRequest(ctx, http.MethodPost, "/v1/mccy/internaltransfers", payload)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusAccepted)
	return c.do(ctx, req)
}

// MCCYIntermediaryAgent - information about the intermediary/correspondent bank.
type MCCYIntermediaryAgent struct {
	// Information that identifies the intermediary agent as a financial institution.
	// Mandatory only when intermediary agent details are being provided.
	FinancialInstitutionIdentification struct {
		// Intermediary agent’s Business Identifier Code (BIC).
		// No need to specify if you are providing the intermediary agent’s American Bankers Association (ABA) routing number.
		BIC string `json:"bic,omitempty"`

		// Intermediary agent’s American Bankers Association (ABA) routing number.
		// No need to specify if you are providing the intermediary agent’s Business Identifier Code (BIC).
		ABA string `json:"aba,omitempty"`

		// Intermediary agent’s name.
		Name string `json:"name,omitempty"`

		// Information about the intermediary agent’s address.
		// Mandatory only when financial institution identification details are being provided.
		AddressDetails struct {
			AddressLine1 string `json:"addressLine1,omitempty"`
			AddressLine2 string `json:"addressLine2,omitempty"`
			AddressLine3 string `json:"addressLine3,omitempty"`
			PostCode     string `json:"postCode,omitempty"`
			Country      string `json:"country"`
		} `json:"addressDetails"`
	} `json:"FinancialInstitutionIdentification"`
}

// MCCYPurpose - information about the purpose of the transaction.
type MCCYPurpose struct {
	// Purpose of the transaction in coded form, as specified in the Purpose11Code code set required by ISO 20022 Message Schemas.
	// No need to specify if you are providing this name in free-form text (i.e., proprietary).
	Code string `json:"code"`
	// Purpose of the transaction in free-form text.
	// No need to specify if you are providing this name in coded form (i.e., code).
	Proprietary string `json:"proprietary"`
}

// MCCYRemittanceInformation - remittance information to accompany the payment.
type MCCYRemittanceInformation struct {
	// Additional remittance information in free-form text, to compliment the structured remittance information.
	AdditionalRemittanceInformation string `json:"additionalRemittanceInformation"`
}

type MCCYTransactionPayload struct {
	// Unique identifier provided to ClearBank for each payment.
	EndToEndID string `json:"endToEndId"`

	// Reference provided by the debtor for the payment.
	Reference string `json:"reference"`

	// Instructed payment amount.
	Amount float64 `json:"amount"`

	// The name used to identify the legal owner of the account from which the funds will be debited.
	DebtorName string `json:"debtorName"`

	// Information about the debtor’s address.
	DebtorAddress struct {
		AddressLine1 string `json:"addressLine1"`
		AddressLine2 string `json:"addressLine2"`
		AddressLine3 string `json:"addressLine3,omitempty"`
		PostCode     string `json:"postCode"`
		Country      string `json:"country"`
	} `json:"debtorAddress"`

	// Debtor’s Business Identifier Code (BIC).
	DebtorBic string `json:"debtorBic,omitempty"`

	// The unique identifier for the account.
	DebtorAccountIdentifier MCCYAccountIdentifier `json:"accountIdentifier"`

	// Three-letter ISO currency code for the currency supported by the debtor account.
	DebtorAccountCurrency string `json:"debtorAccountCurrency"`

	// Information about the intermediary/correspondent bank.
	IntermediaryAgent *MCCYIntermediaryAgent `json:"intermediaryAgent,omitempty"`

	// Information about the purpose of the transaction.
	Purpose *MCCYPurpose `json:"purpose,omitempty"`

	// Remittance information to accompany the payment.
	RemittanceInformation *MCCYRemittanceInformation `json:"remittanceInformation,omitempty"`

	// Information about the creditor of the transaction.
	Creditor struct {
		// Creditor’s name.
		Name string `json:"name"`

		// Information about the creditor’s address.
		Address struct {
			AddressLine1 string `json:"addressLine1"`
			AddressLine2 string `json:"addressLine2"`
			AddressLine3 string `json:"addressLine3,omitempty"`
			PostCode     string `json:"postCode"`
			Country      string `json:"country"`
		} `json:"address"`

		// Creditor’s International Bank Account Number.
		// Mandatory only if account number is not specified.
		IBAN string `json:"iban,omitempty"`

		// Creditor’s Account Number.
		// Mandatory only if iban is not specified.
		AccountNumber string `json:"accountNumber,omitempty"`
	} `json:"creditor"`

	CreditorAgent struct {
		// Information that identifies the creditor agent as a financial institution.
		FinancialInstitutionIdentification struct {
			// Creditor agent’s Business Identifier Code (BIC).
			BIC string `json:"bic,omitempty"`

			// Creditor agent’s American Bankers Association (ABA) routing number.
			ABA string `json:"aba,omitempty"`

			// Creditor agent’s Clearing System Id Code.
			ClearingSystemIDCode string `json:"clearingSystemIdCode,omitempty"`

			// Creditor agent’s Member Id for the specified clearing system.
			// Mandatory only when Clearing System Id Code has been provided.
			MemberID string `json:"memberId,omitempty"`

			// Creditor agent’s name.
			Name string `json:"name"`

			// Information about the creditor agent’s address.
			AddressDetails struct {
				AddressLine1 string `json:"addressLine1,omitempty"`
				AddressLine2 string `json:"addressLine2,omitempty"`
				AddressLine3 string `json:"addressLine3,omitempty"`
				PostCode     string `json:"postCode,omitempty"`
				Country      string `json:"country"`
			} `json:"addressDetails"`
		} `json:"financialInstitutionIdentification"`
	} `json:"creditorAgent"`
}

type CreateMCCYTransactionsPayload struct {
	// Unique identifier for the batch in which the payment is being submitted.
	BatchID string `json:"batchId,omitempty"`

	// Three-letter ISO currency code for the outbound payment.
	Currency string `json:"currencyCode"`

	// Array of transactions.
	Transactions []MCCYTransactionPayload `json:"transactions"`
}

func (c *client) InitiateMCCYTransactions(ctx context.Context, payload CreateMCCYTransactionsPayload) error {
	req, err := c.newRequest(ctx, http.MethodPost, "/v1/mccy/payments", payload)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusAccepted)
	return c.do(ctx, req)
}

type CreateMCCYInboundPaymentPayload struct {
	// Currency of the instructed payment. This is the three-letter ISO currency code.
	InstructedCurrency string `json:"instructedCurrency"`

	// Instructed payment amount.
	InstructedAmount float64 `json:"instructedAmount"`

	// Reference provided by the ultimate debtor for the payment.
	Reference string `json:"reference"`

	// Unique identifier supplied by you in your API request to simulate the sending bank’s reference.
	EndToEndId string `json:"endToEndId"`

	// Information about the ultimate creditor of the transaction.
	UltimateCreditor struct {
		// Ultimate creditor’s International Bank Account Number.
		IBAN string `json:"iban"`
		// Ultimate Creditor Name.
		Name string `json:"name"`
		// Ultimate creditor’s address.
		Address string `json:"address"`
	} `json:"ultimateCreditor"`

	// Information about the ultimate debtor of the transaction.
	UltimateDebtor struct {
		// Ultimate debtor’s unique identifier value. This can only be an IBAN or BIC.
		Identifier string `json:"identifier"`
		// Ultimate Creditor Name.
		Name string `json:"name"`
		// Ultimate creditor’s address.
		Address string `json:"address"`
	} `json:"ultimateDebtor"`
}

type MCCYInboundPaymentResponse struct {
	TransactionID            uuid.UUID `json:"transactionId"`
	InstructedCurrency       string    `json:"instructedCurrency"`
	InstructedAmount         float64   `json:"instructedAmount"`
	Reference                string    `json:"reference"`
	EndToEndID               string    `json:"endToEndId"`
	UltimateCreditorIBAN     string    `json:"ultimateCreditorIBAN"`
	UltimateCreditorName     string    `json:"ultimateCreditorName"`
	UltimateDebtorIdentifier string    `json:"ultimateDebtorIdentifier"`
	UltimateDebtorName       string    `json:"ultimateDebtorName"`
	TimestampCreated         Time      `json:"timestampCreated"`
	BankRefSearchable        string    `json:"bankRefSearchable"`
}

func (c *client) InitiateMCCYInboundPayment(ctx context.Context, accountUniqueID string, payload CreateMCCYInboundPaymentPayload) (data MCCYInboundPaymentResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/v1/mccy/inboundpayment/%s", accountUniqueID), payload)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusCreated)
	req.DecodeTo(&data)
	return data, c.do(ctx, req)
}

type MCCYTransactionStatus string

const (
	MCCYTransactionStatusClearing  MCCYTransactionStatus = "Clearing"
	MCCYTransactionStatusSettled   MCCYTransactionStatus = "Settled"
	MCCYTransactionStatusCancelled MCCYTransactionStatus = "Cancelled" // nolint: misspell
)

type MCCYTransactionResponse struct {
	TransactionID       uuid.UUID     `json:"transactionId"`
	AccountID           uuid.UUID     `json:"accountId"`
	VirtualAccountID    uuid.NullUUID `json:"virtualAccountId"`
	EndToEndID          string        `json:"endToEndId"`
	Reference           string        `json:"reference"`
	UltimateBeneficiary struct {
		AccountIdentifiers []MCCYAccountIdentifier `json:"accountIdentifiers"`
		PayeeName          string                  `json:"payeeName"`
	} `json:"ultimateBeneficiary"`
	UltimateRemitter struct {
		AccountIdentifiers []MCCYAccountIdentifier `json:"accountIdentifiers"`
		PayerName          string                  `json:"payerName"`
	} `json:"ultimateRemitter"`
	Amount           float64 `json:"amount"`
	Currency         string  `json:"currency"`
	CurrencyExchange struct {
		Rate     float64 `json:"rate"`
		Margin   float64 `json:"margin"`
		Currency string  `json:"currency"`
		Amount   float64 `json:"amount"`
	} `json:"currencyExchangeRate"`
	ActualPaymentMethod    string                `json:"actualPaymentMethod"`
	RequestedPaymentMethod string                `json:"requestedPaymentMethod"`
	Kind                   string                `json:"kind"`
	CreatedAt              Time                  `json:"createdAt"`
	SettledAt              string                `json:"settledAt"`
	ValueAt                string                `json:"valueAt"`
	CancelledAt            string                `json:"cancelledAt"`
	CancellationCode       string                `json:"cancellationCode"`
	Reason                 string                `json:"reason"`
	Status                 MCCYTransactionStatus `json:"status"`
	AdditionalProperties   []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"additionalProperties"`
	Identifiers []struct {
		Scope      string `json:"scope"`
		Name       string `json:"name"`
		Identifier string `json:"identifier"`
	} `json:"identifiers"`
}

func (c *client) FetchMCCYTransaction(ctx context.Context, trxID uuid.UUID) (data MCCYTransactionResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/mccy/v1/Transactions/%s", trxID.String()), nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)
	return data, c.do(ctx, req)
}

type MCCYTransactionsResponse struct {
	Transactions []MCCYTransactionResponse `json:"transactions"`
}

func (c *client) FetchMCCYTransactionsForAccount(ctx context.Context, accountID uuid.UUID, currency string, params FetchTransactionsParams) (data MCCYTransactionsResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/mccy/v1/Accounts/%s/transactions", accountID.String()), nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	params.ApplyFor(req)
	req.AddQueryParam("currency", currency)
	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)
	return data, c.do(ctx, req)
}

func (c *client) FetchMCCYTransactionsForVirtualAccount(ctx context.Context, virtualAccountID uuid.UUID, currency string, params FetchTransactionsParams) (data MCCYTransactionsResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/mccy/v1/VirtualAccounts/%s/transactions", virtualAccountID.String()), nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	params.ApplyFor(req)
	req.AddQueryParam("currency", currency)
	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)
	return data, c.do(ctx, req)
}
