package clearbank_test

import (
	"testing"

	"github.com/brokeyourbike/clearbank-api-client-go"
	"github.com/stretchr/testify/assert"
)

func TestMCCYAccountBalancesResponse_UnmarshalJSON(t *testing.T) {
	raw := []byte(`{
		"balances": {
			"USD": {
				"available": 10.00,
				"actual": 9.00
			}
		}
	}`)

	resp := clearbank.MCCYAccountBalancesResponse{}
	err := resp.UnmarshalJSON(raw)
	assert.NoError(t, err)
	assert.Len(t, resp.Balances, 1)

	assert.Equal(t, "USD", resp.Balances[0].Currency)
	assert.Equal(t, 10.00, resp.Balances[0].Available)
	assert.Equal(t, 9.00, resp.Balances[0].Actual)
}
