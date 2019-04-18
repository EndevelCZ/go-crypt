package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/storage"
	"github.com/EndevelCZ/go-crypt/pkg/crypto"
	"github.com/EndevelCZ/go-crypt/pkg/gcp"
	"github.com/googleapis/google-cloud-go-testing/storage/stiface"
	"google.golang.org/api/option"
)

const (
	objectName            = "develop/hello.txt.gpg"
	gcsServiceAccountPath = "/Users/adamplansky/workspace/go/src/github.com/EndevelCZ/go-crypt/secrets/cterminal-developer.json"
	bucketName            = "cterminal-store-wallets"
	encryptionText        = "[GO IZ AWESOME] this is encrypted string for multiple parties!"
)

var (
	pubKeys = []string{
		"/Users/adamplansky/workspace/go/src/github.com/EndevelCZ/go-crypt/secrets/adamplansky.asc",
		"/Users/adamplansky/workspace/go/src/github.com/EndevelCZ/go-crypt/secrets/aliceonprem.asc",
		"/Users/adamplansky/workspace/go/src/github.com/EndevelCZ/go-crypt/secrets/karelonprem.asc",
		"/Users/adamplansky/workspace/go/src/github.com/EndevelCZ/go-crypt/secrets/alicepubtest.asc",
	}
)

func main() {

	// ------ export public keys ------
	pKeys := []io.Reader{}
	for _, p := range pubKeys {
		f, err := os.Open(p)
		if err != nil {
			panic(fmt.Errorf("unable to open file: %s %s", p, err))
		}
		pKeys = append(pKeys, f)
		defer f.Close()
	}
	entities, err := crypto.ReadEntities(pKeys...)
	// ------ export public keys ------

	var buf bytes.Buffer
	src := bytes.NewBufferString(encryptionText)
	err = crypto.Encrypt(entities, nil, src, &buf)
	if err != nil {
		panic(fmt.Errorf("unable to encrypt string: %s ", err))
	}
	fmt.Println("encrypted string OK")

	// ------ create client ------
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(gcsServiceAccountPath))
	if err != nil {
		panic(fmt.Sprintf("unable to create gcs client: %v", err))
	}
	gcs := gcp.NewClient(stiface.AdaptClient(client), bucketName, ctx)
	// ------ create client ------

	// ------ upload ------
	// s := bytes.NewBufferString("hello")
	// s := wBufio
	s := &buf
	err = gcs.UploadObject(s, objectName)
	if err != nil {
		panic(fmt.Sprintf("unable to upload file: %v", err))
	}
	// ------ upload ------

}
