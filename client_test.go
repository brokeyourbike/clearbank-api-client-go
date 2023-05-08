package clearbank

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestIdFromContext(t *testing.T) {
	assert.Equal(t, "", RequestIdFromContext(nil)) //lint:ignore SA1012 this is a test
	assert.Equal(t, "", RequestIdFromContext(context.TODO()))
	assert.Equal(t, "123", RequestIdFromContext(RequestIdContext(context.TODO(), "123")))
}
