package crypto

import (
	"reflect"
	"testing"
)

func TestKmsEncryptSymmetric(t *testing.T) {
	// 	// kmsKey := "projects/cterminal-kms/locations/europe-west3/keyRings/cterminal-store/cryptoKeys/store-test-dev-symetric"
	// 	// plaintext, err := ioutil.ReadFile("README.md")
	// 	// check(err)
	// 	// ciphertext, err := crypto.KmsEncryptSymmetric(kmsKey, plaintext)
	type args struct {
		keyName   string
		plaintext []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			desc: "my first test",
			args: args{keyName: "keyname", plaintext: "plaintext"},
			want: []string{"x-goog-header1:true", "x-goog-header2:0"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := KmsEncryptSymmetric(tt.args.keyName, tt.args.plaintext)
			if (err != nil) != tt.wantErr {
				t.Errorf("KmsEncryptSymmetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KmsEncryptSymmetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKmsDecryptSymmetric(t *testing.T) {
	type args struct {
		keyName    string
		ciphertext []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := KmsDecryptSymmetric(tt.args.keyName, tt.args.ciphertext)
			if (err != nil) != tt.wantErr {
				t.Errorf("KmsDecryptSymmetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KmsDecryptSymmetric() = %v, want %v", got, tt.want)
			}
		})
	}
}
