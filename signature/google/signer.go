package google

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"hash/crc32"

	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var (
	ErrRequestCorrupted  = errors.New("AsymmetricSign: request corrupted in-transit")
	ErrResponseCorrupted = errors.New("AsymmetricSign: response corrupted in-transit")
)

// KeyManagementClient is the interface for the Google Cloud KMS client.
type KeyManagementClient interface {
	AsymmetricSign(context.Context, *kmspb.AsymmetricSignRequest, ...gax.CallOption) (*kmspb.AsymmetricSignResponse, error)
}

// Signer is a Google Cloud KMS signer.
type signer struct {
	client  KeyManagementClient
	keyName string
}

// NewSigner creates a new Google Cloud KMS signer.
func NewSigner(client KeyManagementClient, keyName string) *signer {
	return &signer{client: client, keyName: keyName}
}

// Sign signs the given message using the Google Cloud KMS service.
// It returns the signature and an error if the signing fails.
func (g *signer) Sign(ctx context.Context, message []byte) ([]byte, error) {
	// calculate the digest of the message
	digest := sha256.Sum256(message)

	// compute digest's CRC32C
	digestCRC32C := ComputeCRC32(digest[:])

	// build the signing request
	req := &kmspb.AsymmetricSignRequest{
		Name: g.keyName,
		Digest: &kmspb.Digest{
			Digest: &kmspb.Digest_Sha256{
				Sha256: digest[:],
			},
		},
		DigestCrc32C: wrapperspb.Int64(digestCRC32C),
	}

	result, err := g.client.AsymmetricSign(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to sign digest: %w", err)
	}

	// perform integrity verification on result
	if !result.VerifiedDigestCrc32C {
		return nil, ErrRequestCorrupted
	}

	// perform integrity verification on signature
	if ComputeCRC32(result.Signature) != result.SignatureCrc32C.Value {
		return nil, ErrResponseCorrupted
	}

	encoded := base64.StdEncoding.EncodeToString(result.Signature)

	return []byte(encoded), nil
}

// ComputeCRC32 computes the CRC32C of the given byte slice.
func ComputeCRC32(data []byte) int64 {
	t := crc32.MakeTable(crc32.Castagnoli)

	return int64(crc32.Checksum(data, t))
}
