package local

import "context"

type nilSigner struct{}

// NewNilSigner creates a signer that does nothing.
func NewNilSigner() *nilSigner {
	return &nilSigner{}
}

// Sign returns the same message and no error.
func (s *nilSigner) Sign(ctx context.Context, message []byte) ([]byte, error) {
	return message, nil
}
