package clearbank

import "context"

type Signer interface {
	// Sign signs the given message.
	// It returns the signature and an error if the signing fails.
	Sign(ctx context.Context, message []byte) (signature []byte, err error)
}

type Verifier interface {
	// Verify verifies the given message against signature. It returns an error if the verification fails.
	Verify(ctx context.Context, message []byte, signature []byte) error
}
