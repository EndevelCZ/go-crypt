package crypto

import (
	"context"

	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

func KmsEncryptSymmetric(keyName string, plaintext []byte) ([]byte, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}

	// Build the request.
	req := &kmspb.EncryptRequest{
		Name:      keyName,
		Plaintext: plaintext,
	}
	// Call the API.
	resp, err := client.Encrypt(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Ciphertext, nil
}
func KmsDecryptSymmetric(keyName string, ciphertext []byte) ([]byte, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}

	// Build the request.
	req := &kmspb.DecryptRequest{
		Name:       keyName,
		Ciphertext: ciphertext,
	}
	// Call the API.
	resp, err := client.Decrypt(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Plaintext, nil
}
