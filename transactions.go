package clearbank

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type TransactionsClient interface {
	InitiateFPSTransactions(context.Context, CreateTransactionsPayload) (TransactionsCreatedResponse, error)
	InitiateCHAPSTransactions(context.Context, CreateTransactionsPayload) (TransactionsCreatedResponse, error)
	InitiateFPSPaymentOriginatedOverseas(context.Context, CreateFPSPaymentOriginatedOverseasPayload) (TransactionsCreatedResponse, error)
	FetchTransactions(ctx context.Context, params FetchTransactionsParams) (TransactionsResponse, error)
	FetchTransactionsForAccount(ctx context.Context, accountID uuid.UUID, params FetchTransactionsParams) (TransactionsResponse, error)
	FetchTransactionsForVirtualAccount(ctx context.Context, accountID, virtualAccountID uuid.UUID, params FetchTransactionsParams) (TransactionsResponse, error)
}

type CreditTransfer struct {
	// Identification of the payment instruction.
	PaymentIdentification struct {
		// Unique identification, as assigned by an instructing party
		// for an instructed party, to unambiguously identify the instruction.
		InstructionIdentification string `json:"instructionIdentification"`

		// Unique identification, as assigned by the initiating party to unambiguously identify the transaction.
		// This identification is passed on unchanged throughout the entire end-to-end chain.
		EndToEndIdentification string `json:"endToEndIdentification"`
	} `json:"paymentIdentification"`

	// Indicates the amount and the currency used in a given transaction or account balance.
	Amount struct {
		Currency         string  `json:"currency"`
		InstructedAmount float64 `json:"instructedAmount"`
	} `json:"amount"`

	// The name and, optionally, the legal entity identifier of the account.
	Creditor struct {
		// The name of the account holder.
		Name string `json:"name"`

		// The legal entity identifier of the account holder (eg: charity number).
		// This should be supplied if known.
		LegalEntityIndentifier string `json:"legalEntityIdentifier,omitempty"`
	} `json:"creditor"`

	// Information about the counterpart in a given transaction.
	CreditorAccount struct {
		// The identifiable information of an account.
		Identification struct {
			// The International Bank Account Number (IBAN).
			IBAN string `json:"iban,omitempty"`
			// Unique identification of an account,
			// as assigned by the account servicer,
			// using an identification scheme.
			Other struct {
				// Identification assigned by an institution.
				Identification string `json:"identification"`
				// Name of the identification scheme.
				SchemeName struct {
					// Name of the identification scheme in coded form.
					Code string `json:"code"`
					// Name of the identification scheme in free-form text.
					Proprietary string `json:"proprietary"`
				} `json:"schemeName"`
			} `json:"other,omitempty"`
		} `json:"identification"`
	} `json:"creditorAccount"`

	// Information supplied to enable the matching/reconciliation of an entry
	// with the items that the payment is intended to settle,
	// such as commercial invoices in an accounts’ receivable system.
	RemittanceInformation struct {
		// Information supplied by the remitter (in a structured form),
		// to reconcile an entry with item(s) that the payment intends to settle
		//  (e.g., a purchase reference number).
		Structured struct {
			// Reference information provided by the ultimate debtor
			// to allow the identification of underlying documents by the creditor.
			CreditorReferenceInformation struct {
				// A reference, as assigned by the ultimate debtor
				// to unambiguously refer to the payment transaction.
				// Conditionally required if supplied by the ultimate debtor.
				Reference string `json:"reference"`
			} `json:"creditorReferenceInformation"`
		} `json:"structured,omitempty"`
	} `json:"remittanceInformation,omitempty"`
}

type PaymentInstruction struct {
	// Details about the account holder.
	Debtor struct {
		// The legal entity identifier of the account holder (eg: charity number).
		// This should be supplied if known.
		LegalEntityIdentifier string `json:"legalEntityIdentifier,omitempty"`

		// The name of the account holder.
		Name string `json:"name,omitempty"`

		// The address of the account holder.
		Address string `json:"address"`
	} `json:"debtor"`

	// Information about the counterpart in a given transaction.
	DebtorAccount struct {
		// The identifiable information of an account.
		Identification struct {
			// The International Bank Account Number (IBAN).
			IBAN string `json:"iban"`
		} `json:"identification"`
	} `json:"debtorAccount"`

	// A series of payments that should be made from the debtor account.
	CreditTransfers []CreditTransfer `json:"creditTransfers"`

	// The unique identifier for the payment instruction.
	PaymentInstructionIdentification string `json:"paymentInstructionIdentification"`
}

type CreateTransactionsPayload struct {
	// Details of the payments to be made.
	PaymentInstructions []PaymentInstruction `json:"paymentInstructions"`
}

type TransactionsCreatedResponse struct {
	Transactions []struct {
		EndToEndIdentification string `json:"endToEndIdentification"`
		Response               string `json:"response"`
	}
}

func (c *client) InitiateFPSTransactions(ctx context.Context, payload CreateTransactionsPayload) (data TransactionsCreatedResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/v3/Payments/FPS", payload)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusAccepted)
	req.DecodeTo(&data)

	return data, c.do(ctx, req)
}

func (c *client) InitiateCHAPSTransactions(ctx context.Context, payload CreateTransactionsPayload) (data TransactionsCreatedResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/v3/Payments/CHAPS", payload)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusAccepted)
	req.DecodeTo(&data)

	return data, c.do(ctx, req)
}

type POOPaymentInstruction struct {
	// Information about the debtor of the transaction (the counterparty).
	Debtor struct {
		// The debtor's legal entity identifier (e.g., charity number).
		// Should be supplied if known.
		LegalEntityIdentifier string `json:"legalEntityIdentifier,omitempty"`
	} `json:"debtor"`

	// Information about the counterpart in a given transaction.
	DebtorAccount struct {
		// The identifiable information of an account.
		Identification struct {
			// The International Bank Account Number (IBAN).
			IBAN string `json:"iban"`
		} `json:"identification"`
	} `json:"debtorAccount"`

	// Information about the ultimate debtor of the transaction
	// (the party holding the ultimate debtor account).
	UltimateDebtor struct {
		// Ultimate debtor's name.
		Name string `json:"name"`
		// The full postal address of the ultimate debtor.
		// Depending on the resident country convention, this must include the unit number (if applicable),
		// building name or number, street, town, city, state/province/municipality, and postal code/zip code.
		Address string `json:"address"`
	} `json:"ultimateDebtor"`

	// Information about the ultimate debtor's account.
	UltimateDebtorAccount struct {
		// Account identifiers used to uniquely identify the ultimate debtor's account.
		Identification struct {
			// Ultimate debtor's Bank Identifier Code (BIC).
			BIC string `json:"bic"`
			// Ultimate debtor's account number.
			AccountNumber string `json:"accountNumber"`
		} `json:"identification"`
	} `json:"ultimateDebtorAccount"`

	// Information about the payment to be made from the ultimate debtor's account.
	CreditTransfer struct {
		// Identification of the payment instruction.
		PaymentIdentification struct {
			// Unique identification, as assigned by the initiating party to unambiguously identify the transaction.
			// This identification is passed on unchanged throughout the entire end-to-end chain.
			EndToEndIdentification string `json:"endToEndIdentification"`
		} `json:"paymentIdentification"`

		// Information about the amount (original and instructed), currency and exchange rate.
		Amount struct {
			// If supplied, this must not be 'GBP' and 'originalAmount' as well as 'exchangeRate' must be greater than 0.
			// If not supplied, then 'originalAmount' and 'exchangeRate' must be 0.
			OriginalCurrency string `json:"originalCurrency,omitempty"`
			// Amount of funds to be moved between the debtor and creditor, prior to deduction of charges.
			OriginalAmount float64 `json:"originalAmount,omitempty"`
			// Exchange rate applied to convert the currency of the source amount or OriginalCurrency (prior to deduction of charges) to GBP.
			ExchangeRate     float64 `json:"exchangeRate,omitempty"`
			InstructedAmount float64 `json:"instructedAmount"`
		} `json:"amount"`

		// Information about the creditor of the transaction (the beneficiary account holder).
		Creditor struct {
			// Creditor's name.
			Name string `json:"name"`

			// The creditor's legal entity identifier (e.g., charity number).
			// Should be supplied if known.
			LegalEntityIndentifier string `json:"legalEntityIdentifier,omitempty"`

			// Creditor's address.
			Address string `json:"address,omitempty"`
		} `json:"creditor"`

		// Information about the creditor's account.
		CreditorAccount struct {
			// Account identifiers used to uniquely identify the account.
			Identification struct {
				// The International Bank Account Number (IBAN).
				IBAN string `json:"iban,omitempty"`
			} `json:"identification"`
		} `json:"creditorAccount"`

		// Information supplied by the remitter to reconcile an entry with item(s) that the payment intends to settle.
		RemittanceInformation struct {
			// Information supplied by the remitter (in a structured form),
			// to reconcile an entry with item(s) that the payment intends to settle (e.g., a purchase reference number).
			Structured struct {
				// Reference information provided by the ultimate debtor
				// to allow the identification of underlying documents by the creditor.
				CreditorReferenceInformation struct {
					// A reference, as assigned by the ultimate debtor
					// to unambiguously refer to the payment transaction.
					// Conditionally required if supplied by the ultimate debtor.
					Reference string `json:"reference"`
				} `json:"creditorReferenceInformation"`
			} `json:"structured,omitempty"`

			// Additional unstructured remittance information.
			Unstructured struct {
				// Additional remittance information.
				AdditionalReferenceInformation struct {
					// Additional reference information.
					// Conditionally required if supplied by the ultimate debtor.
					Reference string `json:"reference"`
				} `json:"additionalReferenceInformation"`
			} `json:"unstructured,omitempty"`
		} `json:"remittanceInformation,omitempty"`
	} `json:"creditTransfer"`
}

type CreateFPSPaymentOriginatedOverseasPayload struct {
	// Information about the single payment.
	PaymentInstruction POOPaymentInstruction `json:"paymentInstruction"`
}

func (c *client) InitiateFPSPaymentOriginatedOverseas(ctx context.Context, payload CreateFPSPaymentOriginatedOverseasPayload) (data TransactionsCreatedResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/v2/payments/fps/singlepayment", payload)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusAccepted)
	req.DecodeTo(&data)

	return data, c.do(ctx, req)
}

type FetchTransactionsParams struct {
	PageNumber    int
	PageSize      int
	StartDateTime time.Time
	EndDateTime   time.Time
}

// ApplyFor applies values to the given request as query parameters.
func (p FetchTransactionsParams) ApplyFor(req *request) {
	if p.PageNumber != 0 {
		req.AddQueryParam("pageNumber", strconv.Itoa(p.PageNumber))
	}
	if p.PageSize != 0 {
		req.AddQueryParam("pageSize", strconv.Itoa(p.PageSize))
	}
	if !p.StartDateTime.IsZero() {
		req.AddQueryParam("startDateTime", p.StartDateTime.Format("2006-01-02"))
	}
	if !p.EndDateTime.IsZero() {
		req.AddQueryParam("endDateTime", p.EndDateTime.Format("2006-01-02"))
	}
}

type TransactionResponse struct {
	Amount struct {
		InstructedAmount float64 `json:"instructedAmount"`
		Currency         string  `json:"currency"`
	} `json:"amount"`
	CounterpartAccount struct {
		Identification struct {
			IBAN          string `json:"iban"`
			AccountName   string `json:"accountName"`
			SortCode      string `json:"sortCode"`
			AccountNumber string `json:"accountNumber"`
			Reference     string `json:"reference"`
		} `json:"identification"`
	} `json:"counterpartAccount"`
	DebitCreditCode            string    `json:"debitCreditCode"`
	EndToEndIdentifier         string    `json:"endToEndIdentifier"`
	TransactionID              uuid.UUID `json:"transactionId"`
	TransactionReference       string    `json:"transactionReference"`
	TransactionTime            time.Time `json:"transactionTime"`
	Status                     string    `json:"status"`
	UltimateBeneficiaryAccount *struct {
		ID   uuid.UUID `json:"id"`
		IBAN string    `json:"iban"`
	} `json:"ultimateBeneficiaryAccount"`
	UltimateRemitterAccount *struct {
		ID   uuid.UUID `json:"id"`
		IBAN string    `json:"iban"`
	} `json:"ultimateRemitterAccount"`
}

type TransactionsResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
}

func (c *client) FetchTransactions(ctx context.Context, params FetchTransactionsParams) (data TransactionsResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/v2/Transactions", nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	params.ApplyFor(req)
	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)
	return data, c.do(ctx, req)
}

func (c *client) FetchTransactionsForAccount(ctx context.Context, accountID uuid.UUID, params FetchTransactionsParams) (data TransactionsResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v2/Accounts/%s/Transactions", accountID.String()), nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	params.ApplyFor(req)
	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)
	return data, c.do(ctx, req)
}

func (c *client) FetchTransactionsForVirtualAccount(ctx context.Context, accountID, virtualAccountID uuid.UUID, params FetchTransactionsParams) (data TransactionsResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v2/Accounts/%s/Virtual/%s/Transactions", accountID.String(), virtualAccountID.String()), nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	params.ApplyFor(req)
	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)
	return data, c.do(ctx, req)
}
