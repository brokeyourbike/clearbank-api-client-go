package webhook

// MCCYInstitutionStatementPayload
// This webhook allows you to download a statement for all accounts associated with your institution
type MCCYInstitutionStatementPayload struct {
	URI            string `json:"Uri" validate:"required"`
	ValidUntilDate string `json:"ValidUntilDate" validate:"required"`
	Year           int    `json:"Year" validate:"required"`
	Month          int    `json:"Month" validate:"required"`
	Currency       string `json:"Currency" validate:"required"`
	Format         string `json:"Format" validate:"required"`
}

// MCCYAccountStatementPayload
// This webhook allows you to download a statement for a specific account associated with your institution
type MCCYAccountStatementPayload struct {
	URI            string `json:"Uri" validate:"required"`
	ValidUntilDate string `json:"ValidUntilDate" validate:"required"`
	Year           int    `json:"Year" validate:"required"`
	Month          int    `json:"Month" validate:"required"`
	Currency       string `json:"Currency" validate:"required"`
	Format         string `json:"Format" validate:"required"`
	AccountIBAN    string `json:"AccountIban" validate:"required"`
}
