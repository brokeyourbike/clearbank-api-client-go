package clearbank

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type BacsClient interface {
	FetchAccountMandates(ctx context.Context, accountID uuid.UUID, pageNum int, pageSize int) (MandatesResponse, error)
	CancelAccountMandate(ctx context.Context, accountID, mandateID uuid.UUID, reason string) error

	FetchVirtualAccountMandates(ctx context.Context, accountID, virtualAccountID uuid.UUID, pageNum int, pageSize int) (MandatesResponse, error)
}

type MandatesResponse struct {
	DirectDebitMandates []struct {
		MandateID          uuid.UUID `json:"mandateId"`
		PayerName          string    `json:"payerName"`
		PayerBBAN          string    `json:"payerBban"`
		PayerAccountNumber string    `json:"payerAccountNumber"`
		PayerSortCode      string    `json:"payerSortCode"`
		Reference          string    `json:"reference"`
		ServiceUserNumber  string    `json:"serviceUserNumber"`
		OriginatorName     string    `json:"originatorName"`
		MandateType        string    `json:"mandateType"`
		State              string    `json:"state"`
	} `json:"directDebitMandates"`
}

func (c *client) FetchAccountMandates(ctx context.Context, accountID uuid.UUID, pageNum int, pageSize int) (data MandatesResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v2/Accounts/%s/Mandates", accountID), nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.AddQueryParam("pageNumber", strconv.Itoa(pageNum))
	req.AddQueryParam("pageSize", strconv.Itoa(pageSize))
	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)

	return data, c.do(ctx, req)
}

type cancelMandatePayload struct {
	ReasonCode string `json:"reasonCode"`
}

func (c *client) CancelAccountMandate(ctx context.Context, accountID, mandateID uuid.UUID, reason string) error {
	payload := cancelMandatePayload{ReasonCode: reason}
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/v2/Accounts/%s/Mandates/%s", accountID, mandateID), payload)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusAccepted)
	return c.do(ctx, req)
}

func (c *client) FetchVirtualAccountMandates(ctx context.Context, accountID, virtualAccountID uuid.UUID, pageNum int, pageSize int) (data MandatesResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v2/Accounts/%s/Virtual/%s/Mandates", accountID, virtualAccountID), nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.AddQueryParam("pageNumber", strconv.Itoa(pageNum))
	req.AddQueryParam("pageSize", strconv.Itoa(pageSize))
	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)

	return data, c.do(ctx, req)
}
