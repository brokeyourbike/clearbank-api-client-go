package clearbank

import (
	"context"
	"fmt"
	"net/http"
)

type TransactionsClient interface {
	InitiateFPSTransactions(context.Context, CreateTransactionsPayload) (TransactionsCreatedResponse, error)
	InitiateCHAPSTransactions(context.Context, CreateTransactionsPayload) (TransactionsCreatedResponse, error)
	InitiateFPSPaymentOriginatedOverseas(context.Context, CreateFPSPaymentOriginatedOverseasPayload) (TransactionsCreatedResponse, error)
	InitiateMCCYTransactions(context.Context, CreateMCCYTransactionsPayload) error
	InitiateInternalTransaction(context.Context, CreateInternalTransactionPayload) error
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

type MCCYAccountIdentifier struct {
	Identifier string `json:"identifier"`
	Kind       string `json:"kind"`
}

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
