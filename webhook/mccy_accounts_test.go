package webhook_test

import (
	"testing"

	"github.com/brokeyourbike/clearbank-api-client-go/webhook"
	"github.com/stretchr/testify/assert"
)

func TestMCCYAccountCreatedPayload_GetIdentifier(t *testing.T) {
	p1 := webhook.MCCYAccountCreatedPayload{}
	assert.Equal(t, "", p1.GetIdentifier())

	p2 := webhook.MCCYAccountCreatedPayload{Identifiers: []webhook.MCCYIdentifier{{Identifier: "A"}, {Identifier: "B"}}}
	assert.Equal(t, "A", p2.GetIdentifier())
}

func TestMCCYAccountUpdatedPayload_GetIdentifier(t *testing.T) {
	p1 := webhook.MCCYAccountUpdatedPayload{}
	assert.Equal(t, "", p1.GetIdentifier())

	p2 := webhook.MCCYAccountUpdatedPayload{Identifiers: []webhook.MCCYIdentifier{{Identifier: "A"}, {Identifier: "B"}}}
	assert.Equal(t, "A", p2.GetIdentifier())
}

func TestMCCYVirtualAccountCreatedPayload_GetIdentifier(t *testing.T) {
	p1 := webhook.MCCYVirtualAccountCreatedPayload{}
	assert.Equal(t, "", p1.GetIdentifier())

	p2 := webhook.MCCYVirtualAccountCreatedPayload{Identifiers: []webhook.MCCYIdentifier{{Identifier: "A"}, {Identifier: "B"}}}
	assert.Equal(t, "A", p2.GetIdentifier())
}

func TestMCCYVirtualAccountUpdatedPayload_GetIdentifier(t *testing.T) {
	p1 := webhook.MCCYVirtualAccountUpdatedPayload{}
	assert.Equal(t, "", p1.GetIdentifier())

	p2 := webhook.MCCYVirtualAccountUpdatedPayload{Identifiers: []webhook.MCCYIdentifier{{Identifier: "A"}, {Identifier: "B"}}}
	assert.Equal(t, "A", p2.GetIdentifier())
}
