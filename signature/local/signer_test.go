package local_test

import (
	"context"
	"testing"

	"github.com/brokeyourbike/clearbank-api-client-go/signature/local"
	"github.com/stretchr/testify/assert"
)

func TestNilSigner(t *testing.T) {
	signer := local.NewNilSigner()
	got, err := signer.Sign(context.TODO(), []byte("hello!"))

	assert.NoError(t, err)
	assert.Equal(t, []byte("hello!"), got)
}
