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
	InitiateCHAPSPayment(context.Context, CreateCHAPSPaymentPayload) (PaymentInitiatedResponse, error)
	InitiateFPSPaymentOriginatedOverseas(context.Context, CreateFPSPaymentOriginatedOverseasPayload) (TransactionsCreatedResponse, error)
	FetchTransactions(ctx context.Context, params FetchTransactionsParams) (TransactionsResponse, error)
	FetchTransactionForAccount(ctx context.Context, accountID uuid.UUID, trxID uuid.UUID) (TransactionResponse, error)
	FetchTransactionsForAccount(ctx context.Context, accountID uuid.UUID, params FetchTransactionsParams) (TransactionsResponse, error)
	FetchTransactionForVirtualAccount(ctx context.Context, accountID, virtualAccountID uuid.UUID, trxID uuid.UUID) (TransactionResponse, error)
	FetchTransactionsForVirtualAccount(ctx context.Context, accountID, virtualAccountID uuid.UUID, params FetchTransactionsParams) (TransactionsResponse, error)
}

type CreditTransferOtherIdentification struct {
	// Identification assigned by an institution.
	Identification string `json:"identification"`
	// Name of the identification scheme.
	SchemeName struct {
		// Name of the identification scheme in coded form.
		Code string `json:"code"`
		// Name of the identification scheme in free-form text.
		Proprietary string `json:"proprietary"`
	} `json:"schemeName"`
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
			Other *CreditTransferOtherIdentification `json:"other,omitempty"`
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

const SchemeNameOther = "Other"
const ProprietarySortCode = "SortcodeAccountNumber"

type Address struct {
	BuildingNumber string `json:"buildingNumber,omitempty"`
	BuildingName   string `json:"buildingName,omitempty"`
	StreetName     string `json:"streetName,omitempty"`
	TownName       string `json:"townName"`
	PostCode       string `json:"postCode"`
	Country        string `json:"country"`
}

type CHAPSAccount struct {
	IBAN           string `json:"iban,omitempty"`
	SchemeName     string `json:"schemeName,omitempty"`
	Proprietary    string `json:"proprietary,omitempty"`
	Identification string `json:"identification,omitempty"`
}

type OrganisationIdentifier struct {
	BIC                 string `json:"bic,omitempty"`
	LEI                 string `json:"lei,omitempty"`
	OtherIdentification string `json:"otherIdentification,omitempty"`
}

type PrivateIdentifier struct {
	OtherIdentification string `json:"otherIdentification,omitempty"`
}

type CHAPSParty struct {
	Name                   string                  `json:"name"`
	Address                Address                 `json:"postalAddress"`
	OrganisationIdentifier *OrganisationIdentifier `json:"organisationIdentifier,omitempty"`
	PrivateIdentifier      *PrivateIdentifier      `json:"privateIdentifier,omitempty"`
}

type CreateCHAPSPaymentPayload struct {
	// Unique identification, as assigned by you, to unambiguously identify the payment instruction.
	InstructionIdentification string `json:"instructionIdentification"`

	// Unique identification, as assigned by the initiating party, to unambiguously identify the transaction.
	// This identification is passed on, unchanged, throughout the entire end-to-end chain.
	EndToEndIdentification string `json:"endToEndIdentification"`

	// Amount of money to be moved between the debtor and creditor,
	// before deduction of charges, expressed in the currency as ordered by the initiating party.
	InstructedAmount struct {
		Currency string  `json:"currency"`
		Amount   float64 `json:"amount"`
	} `json:"instructedAmount"`

	// The ClearBank account that will be credited or debited
	// based on the successful completion of the payment instruction.
	// You need to include EITHER the iban field OR the schemeName and identification fields in this object.
	SourceAccount CHAPSAccount `json:"sourceAccount"`

	// Unambiguous identification of the account of the debtor
	// to which a debit entry will be made as a result of the transaction.
	// You need to include EITHER the iban field OR the schemeName and identification fields in this object.
	DebtorAccount CHAPSAccount `json:"debtorAccount"`

	// Unambiguous identification of the account of the creditor
	// to which a credit entry will be posted as a result of the payment transaction.
	CreditorAccount CHAPSAccount `json:"creditorAccount"`

	// Party that is owed an amount of money by the (ultimate) debtor.
	Creditor CHAPSParty `json:"creditor"`

	// Party that owes an amount of money to the (ultimate) creditor.
	Debtor CHAPSParty `json:"debtor"`

	// Underlying reason for the payment transaction, as published in an external purpose code list.
	Purpose string `json:"purpose,omitempty"`

	// Category purpose, in a proprietary form.
	CategoryPurpose string `json:"categoryPurpose,omitempty"`

	// Information supplied to enable the matching of an entry with the items that the transfer is intended to settle,
	// such as commercial invoices in an accounts' receivable system.
	RemittanceInformation struct {
		CreditorReferenceInformation string `json:"creditorReferenceInformation"`
	} `json:"remittanceInformation"`
}

type PaymentInitiatedResponse struct {
	PaymentID uuid.UUID `json:"paymentId"`
}

func (c *client) InitiateCHAPSPayment(ctx context.Context, payload CreateCHAPSPaymentPayload) (data PaymentInitiatedResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/payments/chaps/v5/customer-payments", payload)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusAccepted)
	req.DecodeTo(&data)
	return data, c.do(ctx, req)
}

type ReturnCHAPSPaymentPayload struct {
	PaymentID                 string `json:"paymentId"`
	Reason                    string `json:"reason"`
	InstructionIdentification string `json:"instructionIdentification"`
}

func (c *client) ReturnCHAPSPayment(ctx context.Context, payload ReturnCHAPSPaymentPayload) (data PaymentInitiatedResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/payments/chaps/v5/return-payments", payload)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusAccepted)
	req.DecodeTo(&data)
	return data, c.do(ctx, req)
}

type SimulateCHAPSPaymentPayload struct {
	// Unique identification, as assigned by an instructing party for an instructed party, to unambiguously identify the instruction.
	InstructionIdentification string `json:"instructionIdentification"`

	// Unique identification, as assigned by the initiating party, to unambiguously identify the transaction.
	// This identification is passed on, unchanged, throughout the entire end-to-end chain.
	EndToEndIdentification string `json:"endToEndIdentification"`

	// Amount of money to be moved between the debtor and creditor,
	// before deduction of charges, expressed in the currency as ordered by the initiating party.
	InstructedAmount struct {
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
	} `json:"instructedAmount"`

	// Valid BIC for the debtor account.
	DebtorBIC string `json:"debtorBic"`

	// Unambiguous identification of the account of the debtor to which a debit entry will be made as a result of the transaction.
	// You need to include EITHER the iban field OR the schemeName and identification fields in this object.
	DebtorAccount CHAPSAccount `json:"debtorAccount"`

	// Party that owes an amount of money to the (ultimate) creditor.
	Debtor CHAPSParty `json:"debtor"`

	// Valid BIC for the creditor account.
	CreditorBIC string `json:"creditorBic"`

	// Unambiguous identification of the account of the creditor to which a credit entry will be made as a result of the transaction.
	// You need to include EITHER the iban field OR the schemeName and identification fields in this object.
	CreditorAccount CHAPSAccount `json:"creditorAccount"`

	// Party that is owed an amount of money by the (ultimate) debtor.
	Creditor CHAPSParty `json:"creditor"`

	// Underlying reason for the payment transaction, as published in an external purpose code list.
	Purpose string `json:"purpose"`

	// Category purpose, in a proprietary form.
	CategoryPurpose string `json:"categoryPurpose"`

	// Information supplied to enable the matching of an entry with the items that the transfer is intended to settle,
	// such as commercial invoices in an accounts' receivable system.
	RemittanceInformation struct {
		CreditorReferenceInformation string `json:"creditorReferenceInformation"`
	} `json:"remittanceInformation"`
}

type SimulatedCHAPSPaymentResponse struct {
	UETR string `json:"uetr"`
}

func (c *client) SimulateCHAPSPayment(ctx context.Context, payload SimulateCHAPSPaymentPayload) (data SimulatedCHAPSPaymentResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/inbound-payment-simulation/chaps/v1/customer-payments", payload)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusAccepted)
	req.DecodeTo(&data)
	return data, c.do(ctx, req)
}

type CreateFPSPaymentOriginatedOverseasPayload struct {
	// Unique identification, as assigned by you, to unambiguously identify the payment instruction.
	InstructionIdentification string `json:"instructionIdentification"`

	// Unique identification, as assigned by the initiating party, to unambiguously identify the transaction. This identification is passed on, unchanged, throughout the entire end-to-end chain.
	EndToEndIdentification string `json:"endToEndIdentification"`

	// Underlying reason for the payment transaction, as published in an external purpose code list.
	Purpose string `json:"purpose,omitempty"`

	// Category purpose, in a proprietary form.
	CategoryPurpose string `json:"categoryPurpose,omitempty"`

	// Amount of money to be moved between the debtor and creditor,
	// before deduction of charges, expressed in the currency as ordered by the initiating party.
	InstructedAmount struct {
		Currency string  `json:"currency"`
		Amount   float64 `json:"amount"`
	} `json:"instructedAmount"`

	// The ClearBank account that will be credited or debited
	// based on the successful completion of the payment instruction.
	// You need to include EITHER the iban field OR the schemeName and identification fields in this object.
	SourceAccount struct {
		IBAN           string `json:"iban,omitempty"`
		SchemeName     string `json:"schemeName,omitempty"`
		Proprietary    string `json:"proprietary,omitempty"`
		Identification string `json:"identification,omitempty"`
	} `json:"sourceAccount"`

	// Unambiguous identification of the account of the debtor
	// to which a debit entry will be made as a result of the transaction.
	// You need to include EITHER the iban field OR the schemeName and identification fields in this object.
	DebtorAccount struct {
		IBAN           string `json:"iban,omitempty"`
		SchemeName     string `json:"schemeName,omitempty"`
		Proprietary    string `json:"proprietary,omitempty"`
		Identification string `json:"identification,omitempty"`
	} `json:"debtorAccount"`

	// Unambiguous identification of the account of the creditor
	// to which a credit entry will be posted as a result of the payment transaction.
	CreditorAccount struct {
		IBAN           string `json:"iban,omitempty"`
		SchemeName     string `json:"schemeName,omitempty"`
		Proprietary    string `json:"proprietary,omitempty"`
		Identification string `json:"identification,omitempty"`
	} `json:"creditorAccount"`

	// Party that owes an amount of money to the (ultimate) creditor.
	Debtor struct {
		Name    string  `json:"name"`
		Address Address `json:"postalAddress"`
	} `json:"debtor"`

	// Party that is owed an amount of money by the (ultimate) debtor.
	Creditor struct {
		Name    string  `json:"name"`
		BIC     string  `json:"bic"`
		Address Address `json:"postalAddress"`
	} `json:"creditor"`

	RemittanceInformation struct {
		CreditorReferenceInformation string `json:"creditorReferenceInformation"`
	} `json:"remittanceInformation"`

	// Agent between the debtor's agent and the creditor's agent.
	// You must include the Direct Participant's information in this object.
	IntermediaryAgent struct {
		Name    string  `json:"name"`
		BIC     string  `json:"bic"`
		Address Address `json:"postalAddress"`
	} `json:"intermediaryAgent1"`
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

type TransactionResponseData struct {
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

func (r TransactionResponseData) GetAmount() float64 {
	if r.DebitCreditCode == "CRDT" {
		return r.Amount.InstructedAmount
	}
	return -r.Amount.InstructedAmount
}

type TransactionResponse struct {
	Transaction TransactionResponseData `json:"transaction"`
}

type TransactionsResponse struct {
	Transactions []TransactionResponseData `json:"transactions"`
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

func (c *client) FetchTransactionForAccount(ctx context.Context, accountID uuid.UUID, trxID uuid.UUID) (data TransactionResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v2/Accounts/%s/Transactions/%s", accountID, trxID), nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)
	return data, c.do(ctx, req)
}

func (c *client) FetchTransactionsForAccount(ctx context.Context, accountID uuid.UUID, params FetchTransactionsParams) (data TransactionsResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v2/Accounts/%s/Transactions", accountID), nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	params.ApplyFor(req)
	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)
	return data, c.do(ctx, req)
}

func (c *client) FetchTransactionsForVirtualAccount(ctx context.Context, accountID, virtualAccountID uuid.UUID, params FetchTransactionsParams) (data TransactionsResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v2/Accounts/%s/Virtual/%s/Transactions", accountID, virtualAccountID), nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	params.ApplyFor(req)
	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)
	return data, c.do(ctx, req)
}

func (c *client) FetchTransactionForVirtualAccount(ctx context.Context, accountID, virtualAccountID uuid.UUID, trxID uuid.UUID) (data TransactionResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v2/Accounts/%s/Virtual/%s/Transactions/%s", accountID, virtualAccountID, trxID), nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)
	return data, c.do(ctx, req)
}
