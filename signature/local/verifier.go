package local

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

var (
	ErrFailedToDecodeKey = errors.New("failed to decode public key")
	ErrKeyIsNotPublic    = errors.New("key is not a public key")
	ErrKeyIsNotPublicRSA = errors.New("key is not an RSA public key")
)

type verifier struct {
	key *rsa.PublicKey
}

// NewVerifier creates a new verifier with local public key.
func NewVerifier(key *rsa.PublicKey) verifier {
	return verifier{key: key}
}

// Verify verifies the given message against signature using the local public key.
// Signature must be a base64 encoded string. It returns an error if the verification fails.
func (l verifier) Verify(ctx context.Context, message, signature []byte) error {
	decoded, err := base64.StdEncoding.DecodeString(string(signature))
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}

	// calculate the digest of the message
	digest := sha256.Sum256(message)

	if err := rsa.VerifyPKCS1v15(l.key, crypto.SHA256, digest[:], decoded); err != nil {
		return fmt.Errorf("failed to verify signature: %w", err)
	}

	return nil
}

// ParsePublicKey parses a public RSA key from a PEM encoded bytes.
func ParsePublicKey(pub []byte) (*rsa.PublicKey, error) {
	pubPem, _ := pem.Decode(pub)
	if pubPem == nil {
		return nil, ErrFailedToDecodeKey
	}

	if pubPem.Type != "PUBLIC KEY" {
		return nil, ErrKeyIsNotPublic
	}

	parsed, err := x509.ParsePKIXPublicKey(pubPem.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	key, ok := parsed.(*rsa.PublicKey)
	if !ok {
		return nil, ErrKeyIsNotPublicRSA
	}

	return key, nil
}
