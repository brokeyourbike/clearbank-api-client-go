package clearbank_test

import (
	"testing"

	"github.com/brokeyourbike/clearbank-api-client-go"
	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"TransactionSettled", "2023-07-21T12:07:11.317Z", false},
		{"Fx.Trade.Executed", "2023-07-21T12:07:08Z", false},
		{"Fx.Trade.Executed", "2023-07-21T12:07:11.0036686Z", false},
		{"Payments.Mccy.TransactionCreated", "2023-07-17T07:30:00.65Z", false},
		{"Payments.Mccy.TransactionCreated", "2021-05-22T00:00:00", false},
		{"Payments.Mccy.TransactionSettled", "2023-07-17T07:30:09.1933333Z", false},
		{"Payments.Mccy.TransactionSettled", "2020-09-13T012:22:17.4", false},
		{"Payments.Mccy.TransactionCancelled", "2020-09-13T09:26:15.1762781", false},
		{"Accounts.AccountCreated", "2021-04-08T12:00:00.0000000", false},
		{"Accounts.AccountUpdated", "2020-10-22T20:50:17.7107385", false},
		{"Accounts.VirtualAccountUpdated", "2021-07-21T09:06:33.0628794Z", false},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var d clearbank.Time

			err := d.UnmarshalJSON([]byte(test.value))
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
