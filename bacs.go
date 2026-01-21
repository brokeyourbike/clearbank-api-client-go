package clearbank

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

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

type BacsClient interface {
	FetchAccountMandates(ctx context.Context, accountID uuid.UUID, pageNum int, pageSize int) (MandatesResponse, error)
	FetchVirtualAccountMandates(ctx context.Context, accountID, virtualAccountID uuid.UUID, pageNum int, pageSize int) (MandatesResponse, error)
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
