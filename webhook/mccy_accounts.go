package webhook

import (
	"github.com/brokeyourbike/clearbank-api-client-go"
	"github.com/google/uuid"
)

type IdentifierKind string

const (
	IdentifierKindIBAN       IdentifierKind = "IBAN"
	IdentifierKindBBAN       IdentifierKind = "IBAN"
	IdentifierKindAccountID  IdentifierKind = "AccountId"
	IdentifierKindDescriptor IdentifierKind = "Descriptor"
)

type MCCYIdentifier struct {
	Identifier string         `json:"Identifier" validate:"required"`
	Kind       IdentifierKind `json:"Kind" validate:"required"`
}

// MCCYAccountCreatedPayload
// This webhook confirms that the multicurrency account has been created
type MCCYAccountCreatedPayload struct {
	AccountID        uuid.UUID                   `json:"AccountId" validate:"required"`
	Name             string                      `json:"Name" validate:"required"`
	Label            string                      `json:"Label" validate:"required"`
	Owner            string                      `json:"Owner" validate:"required"`
	Kind             clearbank.MCCYAccountKind   `json:"Kind" validate:"required"`
	Currencies       []string                    `json:"Currencies" validate:"required"`
	Status           clearbank.MCCYAccountStatus `json:"Status" validate:"required"`
	ProductID        string                      `json:"ProductId"`
	CustomerID       string                      `json:"CustomerId"`
	TimestampCreated clearbank.Time              `json:"TimestampCreated" validate:"required"`
	Identifiers      []MCCYIdentifier            `json:"Identifiers" validate:"required"`
	Type             string                      `json:"Type" validate:"required"`
}

// MCCYAccountUpdatedPayload
// This webhook confirms that the multicurrency account has been updated
type MCCYAccountUpdatedPayload struct {
	AccountID         uuid.UUID                   `json:"AccountId" validate:"required"`
	Name              string                      `json:"Name" validate:"required"`
	Label             string                      `json:"Label" validate:"required"`
	Owner             string                      `json:"Owner" validate:"required"`
	Kind              clearbank.MCCYAccountKind   `json:"Kind" validate:"required"`
	Currencies        []string                    `json:"Currencies" validate:"required"`
	Status            clearbank.MCCYAccountStatus `json:"Status" validate:"required"`
	ProductID         string                      `json:"ProductId"`
	CustomerID        string                      `json:"CustomerId"`
	TimestampCreated  clearbank.Time              `json:"TimestampCreated" validate:"required"`
	TimestampModified clearbank.Time              `json:"TimestampModified" validate:"required"`
	StatusReason      string                      `json:"StatusReason"`
	StatusInformation string                      `json:"StatusInformation"`
	Identifiers       []MCCYIdentifier            `json:"Identifiers" validate:"required"`
	Type              string                      `json:"Type" validate:"required"`
}

// MCCYVirtualAccountCreatedPayload
// This webhook confirms that the multicurrency virtual account has been created
type MCCYVirtualAccountCreatedPayload struct {
	AccountID        uuid.UUID                   `json:"AccountId" validate:"required"`
	VirtualAccountID uuid.UUID                   `json:"VirtualAccountId" validate:"required"`
	BatchID          uuid.UUID                   `json:"BatchId" validate:"required"`
	Owner            string                      `json:"Owner" validate:"required"`
	Status           clearbank.MCCYAccountStatus `json:"Status" validate:"required"`
	Identifiers      []MCCYIdentifier            `json:"Identifiers" validate:"required"`
	TimestampCreated clearbank.Time              `json:"TimestampCreated" validate:"required"`
}

// MCCYVirtualAccountUpdatedPayload
// This webhook confirms that the multicurrency virtual account has been updated
type MCCYVirtualAccountUpdatedPayload struct {
	AccountID         uuid.UUID                   `json:"AccountId" validate:"required"`
	VirtualAccountID  uuid.UUID                   `json:"VirtualAccountId" validate:"required"`
	Owner             string                      `json:"Owner" validate:"required"`
	Status            clearbank.MCCYAccountStatus `json:"Status" validate:"required"`
	Identifiers       []MCCYIdentifier            `json:"Identifiers" validate:"required"`
	TimestampCreated  clearbank.Time              `json:"TimestampCreated" validate:"required"`
	TimestampModified clearbank.Time              `json:"TimestampModified" validate:"required"`
	StatusReason      string                      `json:"StatusReason"`
	StatusInformation string                      `json:"StatusInformation"`
}
