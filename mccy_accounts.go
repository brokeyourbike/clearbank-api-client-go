package clearbank

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type MCCYAccountsClient interface {
	// accounts
	FetchMCCYAccount(ctx context.Context, accountID uuid.UUID) (MCCYAccountResponse, error)
	FetchMCCYAccounts(ctx context.Context, pageNum int, pageSize int) (MCCYAccountsResponse, error)
	CreateMCCYAccount(ctx context.Context, payload CreateMCCYAccountPayload) (MCCYAccountResponse, error)
	UpdateMCCYAccount(ctx context.Context, accountID uuid.UUID, payload UpdateMCCYAccountPayload) (MCCYAccountResponse, error)
	UpdateMCCYAccountStatus(ctx context.Context, accountID uuid.UUID, payload UpdateMCCYAccountStatusPayload) (MCCYAccountResponse, error)
	EnableMCCYAccountCurrency(ctx context.Context, accountID uuid.UUID, payload EnableMCCYAccountCurrencyPayload) error
	CloseMCCYAccount(ctx context.Context, accountID uuid.UUID, reason MCCYAccountCloseReason) error
	FetchMCCYAccountBalances(ctx context.Context, accountID uuid.UUID) (MCCYAccountBalancesResponse, error)

	// virtual accounts
	FetchMCCYVirtualAccount(ctx context.Context, virtualAccountID uuid.UUID) (MCCYVirtualAccountResponse, error)
	FetchMCCYVirtualAccounts(ctx context.Context, pageNum int, pageSize int) (MCCYVirtualAccountsResponse, error)
	CreateMCCYVirtualAccount(ctx context.Context, payload CreateMCCYVirtualAccountPayload) (MCCYVirtualAccountResponse, error)
	UpdateMCCYVirtualAccount(ctx context.Context, virtualAccountID uuid.UUID, payload UpdateMCCYVirtualAccountPayload) (MCCYVirtualAccountResponse, error)
	UpdateMCCYVirtualAccountStatus(ctx context.Context, virtualAccountID uuid.UUID, payload UpdateMCCYVirtualAccountStatusPayload) (MCCYVirtualAccountResponse, error)
	CloseMCCYVirtualAccount(ctx context.Context, virtualAccountID uuid.UUID, reason MCCYAccountCloseReason) error
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

func (c *client) FetchMCCYAccount(ctx context.Context, accountID uuid.UUID) (data MCCYAccountResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v2/Accounts/%s/MC", accountID.String()), nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)

	return data, c.do(ctx, req)
}

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

func (c *client) FetchMCCYAccounts(ctx context.Context, pageNum int, pageSize int) (data MCCYAccountsResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/mccy/v1/Accounts", nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.AddQueryParam("pageNumber", strconv.Itoa(pageNum))
	req.AddQueryParam("pageSize", strconv.Itoa(pageSize))
	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)

	return data, c.do(ctx, req)
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

func (c *client) CreateMCCYAccount(ctx context.Context, payload CreateMCCYAccountPayload) (data MCCYAccountResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/mccy/v1/Accounts", payload)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusCreated)
	req.DecodeTo(&data)

	return data, c.do(ctx, req)
}

type UpdateMCCYAccountPayload struct {
	Label string `json:"label"`
	Owner string `json:"owner"`
}

func (c *client) UpdateMCCYAccount(ctx context.Context, accountID uuid.UUID, payload UpdateMCCYAccountPayload) (data MCCYAccountResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodPatch, fmt.Sprintf("/mccy/v1/Accounts/%s", accountID.String()), payload)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)

	return data, c.do(ctx, req)
}

type UpdateMCCYAccountStatusPayload struct {
	Status       string `json:"status"`
	StatusReason string `json:"statusReason,omitempty"`
	Information  string `json:"information,omitempty"`
}

func (c *client) UpdateMCCYAccountStatus(ctx context.Context, accountID uuid.UUID, payload UpdateMCCYAccountStatusPayload) (data MCCYAccountResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodPatch, fmt.Sprintf("/mccy/v1/Accounts/%s/status", accountID.String()), payload)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)

	return data, c.do(ctx, req)
}

type EnableMCCYAccountCurrencyPayload struct {
	Currency string `json:"currency"`
}

func (c *client) EnableMCCYAccountCurrency(ctx context.Context, accountID uuid.UUID, payload EnableMCCYAccountCurrencyPayload) error {
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/mccy/v1/Accounts/%s/currencies", accountID.String()), payload)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusOK)
	return c.do(ctx, req)
}

type MCCYAccountCloseReason string

const (
	MCCYAccountCloseReasonOther                  MCCYAccountCloseReason = "Other"
	MCCYAccountCloseReasonAccountHolderDeceased  MCCYAccountCloseReason = "AccountHolderDeceased"
	MCCYAccountCloseReasonAccountSwitched        MCCYAccountCloseReason = "AccountSwitched"
	MCCYAccountCloseReasonCompanyNoLongerTrading MCCYAccountCloseReason = "CompanyNoLongerTrading"
	MCCYAccountCloseReasonDuplicateAccount       MCCYAccountCloseReason = "DuplicateAccount"
)

func (c *client) CloseMCCYAccount(ctx context.Context, accountID uuid.UUID, reason MCCYAccountCloseReason) error {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/mccy/v1/Accounts/%s", accountID.String()), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.AddQueryParam("accountCloseReason", string(reason))
	req.ExpectStatus(http.StatusNoContent)
	return c.do(ctx, req)
}

type MCCYAccountBalance struct {
	Currency  string  `json:"currency"`
	Available float64 `json:"available"`
	Actual    float64 `json:"actual"`
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

func (c *client) FetchMCCYAccountBalances(ctx context.Context, accountID uuid.UUID) (data MCCYAccountBalancesResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/mccy/v1/Accounts/%s/balances", accountID.String()), nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)

	return data, c.do(ctx, req)
}

type MCCYVirtualAccountResponse struct {
	ID                uuid.UUID               `json:"id"`
	AccountID         uuid.UUID               `json:"accountId"`
	Owner             string                  `json:"owner"`
	Identifiers       []MCCYAccountIdentifier `json:"identifiers"`
	Status            MCCYAccountStatus       `json:"status"`
	StatusReason      string                  `json:"statusReason"`
	StatusInformation string                  `json:"statusInformation"`
}

func (c *client) FetchMCCYVirtualAccount(ctx context.Context, virtualAccountID uuid.UUID) (data MCCYVirtualAccountResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/mccy/v1/VirtualAccounts/%s", virtualAccountID.String()), nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)

	return data, c.do(ctx, req)
}

type MCCYVirtualAccountsResponse struct {
	VirtualAccounts []MCCYVirtualAccountResponse `json:"virtualAccounts"`
}

func (c *client) FetchMCCYVirtualAccounts(ctx context.Context, pageNum int, pageSize int) (data MCCYVirtualAccountsResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/mccy/v1/VirtualAccounts", nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.AddQueryParam("pageNumber", strconv.Itoa(pageNum))
	req.AddQueryParam("pageSize", strconv.Itoa(pageSize))
	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)

	return data, c.do(ctx, req)
}

type CreateMCCYVirtualAccountPayload struct {
	AccountID      string `json:"accountId"`
	VirtualAccount struct {
		Owner string `json:"owner"`
	} `json:"virtualAccount"`
}

func (c *client) CreateMCCYVirtualAccount(ctx context.Context, payload CreateMCCYVirtualAccountPayload) (data MCCYVirtualAccountResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/mccy/v1/VirtualAccounts", payload)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusCreated)
	req.DecodeTo(&data)

	return data, c.do(ctx, req)
}

type UpdateMCCYVirtualAccountPayload struct {
	Owner string `json:"owner"`
}

func (c *client) UpdateMCCYVirtualAccount(ctx context.Context, virtualAccountID uuid.UUID, payload UpdateMCCYVirtualAccountPayload) (data MCCYVirtualAccountResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodPatch, fmt.Sprintf("/mccy/v1/VirtualAccounts/%s", virtualAccountID.String()), payload)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)

	return data, c.do(ctx, req)
}

type UpdateMCCYVirtualAccountStatusPayload struct {
	Status       string `json:"status"`
	StatusReason string `json:"statusReason,omitempty"`
	Information  string `json:"information,omitempty"`
}

func (c *client) UpdateMCCYVirtualAccountStatus(ctx context.Context, virtualAccountID uuid.UUID, payload UpdateMCCYVirtualAccountStatusPayload) (data MCCYVirtualAccountResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodPatch, fmt.Sprintf("/mccy/v1/VirtualAccounts/%s/status", virtualAccountID.String()), payload)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)

	return data, c.do(ctx, req)
}

func (c *client) CloseMCCYVirtualAccount(ctx context.Context, virtualAccountID uuid.UUID, reason MCCYAccountCloseReason) error {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/mccy/v1/VirtualAccounts/%s", virtualAccountID.String()), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.AddQueryParam("accountCloseReason", string(reason))
	req.ExpectStatus(http.StatusNoContent)
	return c.do(ctx, req)
}
