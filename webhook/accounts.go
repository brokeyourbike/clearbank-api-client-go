package webhook

import (
	"github.com/brokeyourbike/clearbank-api-client-go"
	"github.com/google/uuid"
)

type AccountType string

const (
	AccountTypeCurrent          AccountType = "CurrentAccount"
	AccountTypeDeposit          AccountType = "DepositAccount"
	AccountTypeControl          AccountType = "ControlAccount"
	AccountTypeSegregatedClient AccountType = "SegregatedClientAccount"
	AccountTypeInstitution      AccountType = "InstitutionAccount"
)

// AccountCreatedPayload
// This webhook confirms the account has been created
type AccountCreatedPayload struct {
	AccountID          uuid.UUID `json:"AccountId" validate:"required"`
	AccountName        string    `json:"AccountName" validate:"required"`
	AccountHolderLabel string    `json:"AccountHolderLabel" validate:"required"`
	AccountIdentifiers struct {
		IBAN string `json:"IBAN" validate:"required"`
		BBAN string `json:"BBAN" validate:"required"`
	} `json:"AccountIdentifiers" validate:"required"`
	TimestampCreated clearbank.Time `json:"TimestampCreated" validate:"required"`
	AccountType      AccountType    `json:"AccountType" validate:"required"`
}

type DisabledReason string

const (
	DisabledReasonAccountClosed    DisabledReason = "AccountClosed"
	DisabledReasonAccountSuspended DisabledReason = "AccountSuspended"
)

// AccountDisabledPayload
type AccountDisabledPayload struct {
	AccountID         uuid.UUID      `json:"AccountId" validate:"required"`
	DisabledReason    DisabledReason `json:"DisabledReason" validate:"required"`
	DisabledTimestamp clearbank.Time `json:"DisabledTimestamp" validate:"required"`
}

// VirtualAccountCreatedPayload
// This webhook confirms creation of a virtual account.
type VirtualAccountCreatedPayload struct {
	AccountID         uuid.UUID `json:"AccountId" validate:"required"`
	VirtualAccountID  uuid.UUID `json:"VirtualAccountId" validate:"required"`
	AccountIdentifier struct {
		IBAN               string `json:"Iban" validate:"required"`
		BBAN               string `json:"Bban" validate:"required"`
		ExternalIdentifier string `json:"ExternalIdentifier"`
	} `json:"AccountIdentifier" validate:"required"`
	OwnerName        string         `json:"OwnerName"`
	TimestampCreated clearbank.Time `json:"TimestampCreated" validate:"required"`
}

// VirtualAccountCreationFailedPayload
// This webhook confirms a virtual account creation failure
type VirtualAccountCreationFailedPayload struct {
	AccountID         uuid.UUID `json:"AccountId" validate:"required"`
	VirtualAccountID  uuid.UUID `json:"VirtualAccountId" validate:"required"`
	AccountIdentifier struct {
		IBAN               string `json:"Iban" validate:"required"`
		BBAN               string `json:"Bban" validate:"required"`
		ExternalIdentifier string `json:"ExternalIdentifier"`
	} `json:"AccountIdentifier" validate:"required"`
	Errors    map[string]string `json:"Errors" validate:"required"`
	OwnerName string            `json:"OwnerName"`
	Timestamp clearbank.Time    `json:"Timestamp" validate:"required"`
}
