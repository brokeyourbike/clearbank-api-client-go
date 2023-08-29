package clearbank

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type MCCYAccountsClient interface {
	// accounts
	FetchMCCYAccount(context.Context, uuid.UUID) (MCCYAccountResponse, error)
	FetchMCCYAccounts(context.Context, int, int) (MCCYAccountsResponse, error)
	CreateMCCYAccount(context.Context, CreateMCCYAccountPayload) (MCCYAccountResponse, error)
	UpdateMCCYAccount(context.Context, uuid.UUID, UpdateMCCYAccountPayload) (MCCYAccountResponse, error)
	UpdateMCCYAccountStatus(context.Context, uuid.UUID, UpdateMCCYAccountStatusPayload) (MCCYAccountResponse, error)
	EnableMCCYAccountCurrency(context.Context, uuid.UUID, EnableMCCYAccountCurrencyPayload) error
	CloseMCCYAccount(context.Context, uuid.UUID, MCCYAccountCloseReason) error
	FetchMCCYAccountBalances(context.Context, uuid.UUID) (MCCYAccountBalancesResponse, error)

	// virtual accounts
	FetchMCCYVirtualAccount(context.Context, uuid.UUID) (MCCYVirtualAccountResponse, error)
	FetchMCCYVirtualAccounts(context.Context, int, int) (MCCYVirtualAccountsResponse, error)
	CreateMCCYVirtualAccount(context.Context, CreateMCCYVirtualAccountPayload) (MCCYVirtualAccountResponse, error)
	UpdateMCCYVirtualAccount(context.Context, uuid.UUID, UpdateMCCYVirtualAccountPayload) (MCCYVirtualAccountResponse, error)
	UpdateMCCYVirtualAccountStatus(context.Context, uuid.UUID, UpdateMCCYVirtualAccountStatusPayload) (MCCYVirtualAccountResponse, error)
	CloseMCCYVirtualAccount(context.Context, uuid.UUID, MCCYAccountCloseReason) error
}

type MCCYAccountKind string

const (
	MCCYAccountKindYourFunds            MCCYAccountKind = "YourFunds"
	MCCYAccountKindGeneralSegregated    MCCYAccountKind = "GeneralSegregated"
	MCCYAccountKindDesignatedSegregated MCCYAccountKind = "DesignatedSegregated"
	MCCYAccountKindGeneralClient        MCCYAccountKind = "GeneralClient"
	MCCYAccountKindDesignatedClient     MCCYAccountKind = "DesignatedClient"
)

type MCCYAccountStatus string
type MCCYAccountStatusReason string

const (
	MCCYAccountStatusActive                      MCCYAccountStatus       = "Active"
	MCCYAccountStatusSuspended                   MCCYAccountStatus       = "Suspended"
	MCCYAccountStatusClosed                      MCCYAccountStatus       = "Closed"
	MCCYAccountStatusNotProvided                 MCCYAccountStatus       = "NotProvided"
	MCCYAccountStatusReasonAccountHolderBankrupt MCCYAccountStatusReason = "AccountHolderBankrupt"
	MCCYAccountStatusReasonDissatisfiedCustomer  MCCYAccountStatusReason = "DissatisfiedCustomer"
	MCCYAccountStatusReasonFinancialCrime        MCCYAccountStatusReason = "FinancialCrime"
	MCCYAccountStatusReasonFraudFirstParty       MCCYAccountStatusReason = "FraudFirstParty"
	MCCYAccountStatusReasonFraudThirdParty       MCCYAccountStatusReason = "FraudThirdParty"
	MCCYAccountStatusReasonOther                 MCCYAccountStatusReason = "Other"
)

type MCCYAccountIdentifier struct {
	Identifier string `json:"identifier"`
	Kind       string `json:"kind"`
}

type MCCYAccountResponse struct {
	ID                uuid.UUID               `json:"id"`
	Name              string                  `json:"name"`
	Label             string                  `json:"label"`
	Owner             string                  `json:"owner"`
	Kind              MCCYAccountKind         `json:"kind"`
	Currencies        []string                `json:"currencies"`
	Identifiers       []MCCYAccountIdentifier `json:"identifiers"`
	Status            MCCYAccountStatus       `json:"status"`
	StatusReason      string                  `json:"statusReason"`
	StatusInformation string                  `json:"statusInformation"`
}

//

type MCCYAccountsResponse struct {
	Accounts []struct {
		ID                uuid.UUID               `json:"id"`
		Name              string                  `json:"name"`
		Label             string                  `json:"label"`
		Owner             string                  `json:"owner"`
		Kind              MCCYAccountKind         `json:"kind"`
		Currencies        []string                `json:"currencies"`
		Identifiers       []MCCYAccountIdentifier `json:"identifiers"`
		Status            MCCYAccountStatus       `json:"status"`
		StatusReason      string                  `json:"statusReason"`
		StatusInformation string                  `json:"statusInformation"`
		Type              string                  `json:"type"`
	} `json:"accounts"`
}

type RoutingCodeKind string

const RoutingCodeKindGB RoutingCodeKind = "GBSortCode"

type CreateMCCYAccountPayload struct {
	Label       string          `json:"label"`
	Owner       string          `json:"owner"`
	Kind        MCCYAccountKind `json:"kind"`
	Currencies  []string        `json:"currencies"`
	RoutingCode struct {
		Code string          `json:"code"`
		Kind RoutingCodeKind `json:"kind"`
	} `json:"routingCode"`
}

type UpdateMCCYAccountPayload struct {
	Label string `json:"label"`
	Owner string `json:"owner"`
}

type UpdateMCCYAccountStatusPayload struct {
	Status       string `json:"status"`
	StatusReason string `json:"statusReason,omitempty"`
	Information  string `json:"information,omitempty"`
}

type EnableMCCYAccountCurrencyPayload struct {
	Currency string `json:"currency"`
}

//

type MCCYAccountCloseReason string

const (
	MCCYAccountCloseReasonOther                  MCCYAccountCloseReason = "Other"
	MCCYAccountCloseReasonAccountHolderDeceased  MCCYAccountCloseReason = "AccountHolderDeceased"
	MCCYAccountCloseReasonAccountSwitched        MCCYAccountCloseReason = "AccountSwitched"
	MCCYAccountCloseReasonCompanyNoLongerTrading MCCYAccountCloseReason = "CompanyNoLongerTrading"
	MCCYAccountCloseReasonDuplicateAccount       MCCYAccountCloseReason = "DuplicateAccount"
)

//

type MCCYAccountBalance struct {
	Currency  string  `json:"Currency"`
	Available float64 `json:"Available"`
	Actual    float64 `json:"Actual"`
}

type MCCYAccountBalancesResponse struct {
	Balances []MCCYAccountBalance `json:"balances"`
}

func (m *MCCYAccountBalancesResponse) UnmarshalJSON(data []byte) error {
	type payload struct {
		Balances map[string]struct {
			Available float64 `json:"Available"`
			Actual    float64 `json:"Actual"`
		} `json:"balances"`
	}

	var p payload

	if err := json.Unmarshal(data, &p); err != nil {
		return fmt.Errorf("unable to unmarshal MCCY account balances: %w", err)
	}

	for currency, balance := range p.Balances {
		m.Balances = append(m.Balances, MCCYAccountBalance{
			Currency:  currency,
			Available: balance.Available,
			Actual:    balance.Actual,
		})
	}

	return nil
}

//

type MCCYVirtualAccountResponse struct {
	ID                uuid.UUID               `json:"id"`
	AccountID         uuid.UUID               `json:"accountId"`
	Owner             string                  `json:"owner"`
	Identifiers       []MCCYAccountIdentifier `json:"identifiers"`
	Status            MCCYAccountStatus       `json:"status"`
	StatusReason      string                  `json:"statusReason"`
	StatusInformation string                  `json:"statusInformation"`
}

type MCCYVirtualAccountsResponse struct {
	VirtualAccounts []MCCYVirtualAccountResponse `json:"virtualAccounts"`
}

//

type CreateMCCYVirtualAccountPayload struct {
	AccountID      string `json:"accountId"`
	VirtualAccount struct {
		Owner string `json:"owner"`
	} `json:"virtualAccount"`
}

type UpdateMCCYVirtualAccountPayload struct {
	Owner string `json:"owner"`
}

type UpdateMCCYVirtualAccountStatusPayload struct {
	Status       string `json:"status"`
	StatusReason string `json:"statusReason,omitempty"`
	Information  string `json:"information,omitempty"`
}
