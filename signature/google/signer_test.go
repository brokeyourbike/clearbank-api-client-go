package google_test

import (
	"context"
	"errors"
	"testing"

	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/brokeyourbike/clearbank-api-client-go/signature/google"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func Test_Sign(t *testing.T) {
	tests := []struct {
		name      string
		message   []byte
		wantErr   bool
		setupMock func(clientMock *google.MockKeyManagementClient)
	}{
		{"successful signature", []byte("test"), false, func(clientMock *google.MockKeyManagementClient) {
			response := &kmspb.AsymmetricSignResponse{
				VerifiedDigestCrc32C: true,
				Signature:            []byte("signature"),
				SignatureCrc32C:      wrapperspb.Int64(google.ComputeCRC32([]byte("signature"))),
			}

			clientMock.On("AsymmetricSign", mock.Anything, mock.Anything).Return(response, nil).Once()
		}},
		{"failed to sign", []byte("test"), true, func(clientMock *google.MockKeyManagementClient) {
			response := &kmspb.AsymmetricSignResponse{}
			err := errors.New("failed to sign") // nolint: goerr113
			clientMock.On("AsymmetricSign", mock.Anything, mock.Anything).Return(response, err).Once()
		}},
		{"digest verification failed", []byte("test"), true, func(clientMock *google.MockKeyManagementClient) {
			response := &kmspb.AsymmetricSignResponse{
				VerifiedDigestCrc32C: false,
			}
			clientMock.On("AsymmetricSign", mock.Anything, mock.Anything).Return(response, nil).Once()
		}},
		{"signature verification failed", []byte("test"), true, func(clientMock *google.MockKeyManagementClient) {
			response := &kmspb.AsymmetricSignResponse{
				VerifiedDigestCrc32C: true,
				Signature:            []byte("a"),
				SignatureCrc32C:      wrapperspb.Int64(google.ComputeCRC32([]byte("b"))),
			}
			clientMock.On("AsymmetricSign", mock.Anything, mock.Anything).Return(response, nil).Once()
		}},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			clientMock := google.NewMockKeyManagementClient(t)
			test.setupMock(clientMock)

			signer := google.NewSigner(clientMock, "")
			value, err := signer.Sign(context.TODO(), test.message)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, value)
			}
		})
	}
}

func Test_Sign_ReturnsBase64OfSignature(t *testing.T) {
	clientMock := google.NewMockKeyManagementClient(t)
	response := &kmspb.AsymmetricSignResponse{
		VerifiedDigestCrc32C: true,
		Signature:            []byte("signature"),
		SignatureCrc32C:      wrapperspb.Int64(google.ComputeCRC32([]byte("signature"))),
	}
	clientMock.On("AsymmetricSign", mock.Anything, mock.Anything).Return(response, nil).Once()

	signer := google.NewSigner(clientMock, "")
	value, err := signer.Sign(context.TODO(), []byte("test"))

	assert.NoError(t, err)
	assert.Equal(t, []byte("c2lnbmF0dXJl"), value)
}
