package webhook_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/brokeyourbike/clearbank-api-client-go/webhook"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/Mccy.InternalTransfers.Settled.json
var mccy_InternalTransfers_Settled []byte

func TestMCCYInternalTransfersSettledPayload(t *testing.T) {
	var p webhook.MCCYInternalTransfersSettledPayload
	err := json.Unmarshal(mccy_InternalTransfers_Settled, &p)
	require.NoError(t, err)

	assert.Equal(t, "f92af28f-0635-4d81-893f-6c32af1d2a17", p.Reference)
}
