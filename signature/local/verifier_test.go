package local_test

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"sync"
	"testing"

	"github.com/brokeyourbike/clearbank-api-client-go/signature/local"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/invalid.pem
var invalidKeyPem []byte

//go:embed testdata/public-key.pem
var publicKeyPem []byte

//go:embed testdata/public-key-ed25519.pem
var publicKeyEd25519Pem []byte

//go:embed testdata/private-key.pem
var privateKeyPem []byte

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
			verifier := local.NewVerifier(&privateKey.PublicKey)
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

	verifier := local.NewVerifier(&privateKey.PublicKey)

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

func TestParsePublicKey(t *testing.T) {
	key, err := local.ParsePublicKey(publicKeyPem)
	assert.NoError(t, err)
	assert.NotNil(t, key)

	_, err = local.ParsePublicKey(invalidKeyPem)
	assert.Error(t, err)
	assert.ErrorIs(t, err, local.ErrFailedToDecodeKey)

	_, err = local.ParsePublicKey(publicKeyEd25519Pem)
	assert.Error(t, err)
	assert.ErrorIs(t, err, local.ErrKeyIsNotPublicRSA)

	_, err = local.ParsePublicKey(privateKeyPem)
	assert.Error(t, err)
	assert.ErrorIs(t, err, local.ErrKeyIsNotPublic)
}
