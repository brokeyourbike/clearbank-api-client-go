package clearbank

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type AccountsClient interface {
	// accounts
	FetchAccount(ctx context.Context, accountID uuid.UUID) (AccountResponse, error)
	FetchAccounts(ctx context.Context, pageNum int, pageSize int) (AccountsResponse, error)
	CreateAccount(ctx context.Context, payload CreateAccountPayload) (AccountResponse, error)
	UpdateAccount(ctx context.Context, accountID uuid.UUID, payload UpdateAccountPayload) error

	// virtual accounts
	FetchVirtualAccountsFor(ctx context.Context, accountID uuid.UUID, pageNum int, pageSize int) (VirtualAccountsResponse, error)
	CreateVirtualAccounts(ctx context.Context, accountID uuid.UUID, payload CreateVirtualAccountsPayload) error
	UpdateVirtualAccount(ctx context.Context, accountID uuid.UUID, virtualAccountID uuid.UUID, payload UpdateVirtualAccountPayload) error
	DisableVirtualAccount(ctx context.Context, accountID uuid.UUID, virtualAccountID uuid.UUID) error
}

type AccountStatus string
type VirtualAccountStatus string

const (
	AccountStatusNotProvided     AccountStatus        = "NotProvided"
	AccountStatusEnabled         AccountStatus        = "Enabled"
	AccountStatusClosed          AccountStatus        = "Closed"
	AccountStatusSuspended       AccountStatus        = "Suspended"
	VirtualAccountStatusEnabled  VirtualAccountStatus = "Enabled"
	VirtualAccountStatusDisabled VirtualAccountStatus = "Disabled"
)

type AccountBalance struct {
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Status   string  `json:"status"`
}

type AccountResponse struct {
	Account struct {
		ID           uuid.UUID        `json:"id"`
		Name         string           `json:"name"`
		Owner        string           `json:"label"`
		Type         string           `json:"type"`
		Currency     []string         `json:"currency"`
		Balances     []AccountBalance `json:"balances"`
		Status       AccountStatus    `json:"status"`
		StatusReason string           `json:"statusReason"`
		IBAN         string           `json:"iban"`
	} `json:"account"`
}

func (c *client) FetchAccount(ctx context.Context, accountID uuid.UUID) (data AccountResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v3/Accounts/%s", accountID.String()), nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)

	return data, c.do(ctx, req)
}

type AccountsResponse struct {
	Accounts []struct {
		ID           uuid.UUID        `json:"id"`
		Name         string           `json:"name"`
		Owner        string           `json:"label"`
		Type         string           `json:"type"`
		Currency     []string         `json:"currency"`
		Balances     []AccountBalance `json:"balances"`
		Status       AccountStatus    `json:"status"`
		StatusReason string           `json:"statusReason"`
		IBAN         string           `json:"iban"`
	} `json:"accounts"`
}

func (c *client) FetchAccounts(ctx context.Context, pageNum int, pageSize int) (data AccountsResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/v3/Accounts", nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.AddQueryParam("pageNumber", strconv.Itoa(pageNum))
	req.AddQueryParam("pageSize", strconv.Itoa(pageSize))
	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)

	return data, c.do(ctx, req)
}

type AccountUsageType string

const (
	AccountUsageTypeYourFunds            AccountUsageType = "YourFunds"
	AccountUsageTypeSegregatedDesignated AccountUsageType = "SegregatedDesignated"
	AccountUsageTypeSegregatedPooled     AccountUsageType = "SegregatedPooled"
)

type CreateAccountPayload struct {
	Name  string `json:"accountName"`
	Owner struct {
		Name string `json:"name"`
	} `json:"owner"`
	SortCode  string           `json:"sortCode"`
	UsageType AccountUsageType `json:"usageType,omitempty"`
}

func (c *client) CreateAccount(ctx context.Context, payload CreateAccountPayload) (data AccountResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/v3/Accounts", payload)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusCreated)
	req.DecodeTo(&data)

	return data, c.do(ctx, req)
}

type UpdateAccountPayload struct {
	Status           string `json:"status,omitempty"`
	StatusReason     string `json:"statusReason,omitempty"`
	OwnerName        string `json:"ownerName,omitempty"`
	LegalOwnerType   string `json:"legalOwnerType,omitempty"`
	RelationshipType string `json:"relationshipType,omitempty"`
	MinimumBalance   int    `json:"minimumBalance,omitempty"`
}

func (c *client) UpdateAccount(ctx context.Context, accountID uuid.UUID, payload UpdateAccountPayload) error {
	req, err := c.newRequest(ctx, http.MethodPatch, fmt.Sprintf("/v1/Accounts/%s", accountID.String()), payload)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusNoContent)
	return c.do(ctx, req)
}

type VirtualAccountsResponse struct {
	Accounts []struct {
		ID         uuid.UUID            `json:"id"`
		Name       string               `json:"name"`
		Owner      string               `json:"label"`
		Type       string               `json:"type"`
		Status     VirtualAccountStatus `json:"status"`
		Currencies []string             `json:"currency"`
		IBAN       string               `json:"iban"`
	} `json:"accounts"`
}

func (c *client) FetchVirtualAccountsFor(ctx context.Context, accountID uuid.UUID, pageNum int, pageSize int) (data VirtualAccountsResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v2/Accounts/%s/Virtual", accountID.String()), nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.AddQueryParam("pageNumber", strconv.Itoa(pageNum))
	req.AddQueryParam("pageSize", strconv.Itoa(pageSize))
	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)

	return data, c.do(ctx, req)
}

type CreateVirtualAccountPayload struct {
	OwnerName         string `json:"ownerName"`
	AccountIdentifier struct {
		IBAN                  string `json:"iban,omitempty"`
		BBAN                  string `json:"bban,omitempty"`
		ProprietaryIdentifier string `json:"proprietaryIdentifier,omitempty"`
		ExternalIdentifier    string `json:"externalIdentifier,omitempty"`
	} `json:"accountIdentifier"`
}

type CreateVirtualAccountsPayload struct {
	VirtualAccounts []CreateVirtualAccountPayload `json:"virtualAccounts"`
}

func (c *client) CreateVirtualAccounts(ctx context.Context, accountID uuid.UUID, payload CreateVirtualAccountsPayload) error {
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/v2/Accounts/%s/Virtual", accountID.String()), payload)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusAccepted)
	return c.do(ctx, req)
}

type UpdateVirtualAccountPayload struct {
	OwnerName        string `json:"ownerName,omitempty"`
	LegalOwnerType   string `json:"legalOwnerType,omitempty"`
	RelationshipType string `json:"relationshipType,omitempty"`
}

func (c *client) UpdateVirtualAccount(ctx context.Context, accountID uuid.UUID, virtualAccountID uuid.UUID, payload UpdateVirtualAccountPayload) error {
	req, err := c.newRequest(ctx, http.MethodPatch, fmt.Sprintf("/v1/Accounts/%s/Virtual/%s", accountID.String(), virtualAccountID.String()), payload)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusNoContent)
	return c.do(ctx, req)
}

func (c *client) DisableVirtualAccount(ctx context.Context, accountID uuid.UUID, virtualAccountID uuid.UUID) error {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/v1/Accounts/%s/Virtual/%s", accountID.String(), virtualAccountID.String()), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusNoContent)
	return c.do(ctx, req)
}
