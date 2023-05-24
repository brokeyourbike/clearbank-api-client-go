package verifier_test

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"sync"
	"testing"

	"github.com/brokeyourbike/clearbank-api-client-go/verifier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LocalVerifier_Verify(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	tests := []struct {
		name      string
		message   []byte
		signature []byte
		wantErr   bool
	}{
		{"valid signature", []byte("a"), signPKCS1v15(t, privateKey, []byte("a")), false},
		{"invalid signature", []byte("a"), signPKCS1v15(t, privateKey, []byte("b")), true},
		{"signature is not base64", []byte("a"), []byte("..."), true},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			verifier := verifier.NewLocalVerifier(&privateKey.PublicKey)
			err = verifier.Verify(context.TODO(), test.message, test.signature)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_LocalVerifier_Verify_ConcurrenUsage(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	verifier := verifier.NewLocalVerifier(&privateKey.PublicKey)

	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			message := []byte("a")
			sig := signPKCS1v15(t, privateKey, message)

			err := verifier.Verify(context.TODO(), message, sig)
			assert.NoError(t, err)
		}()
	}

	wg.Wait()
}

func signPKCS1v15(t *testing.T, key *rsa.PrivateKey, message []byte) []byte {
	t.Helper()

	hashed := sha256.Sum256(message)
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hashed[:])
	require.NoError(t, err)

	encoded := base64.StdEncoding.EncodeToString(signature)

	return []byte(encoded)
}
