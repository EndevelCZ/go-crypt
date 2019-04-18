package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"cloud.google.com/go/storage"
	"github.com/EndevelCZ/go-crypt/pkg/crypto"
	"github.com/EndevelCZ/go-crypt/pkg/gcp"
	"github.com/googleapis/google-cloud-go-testing/storage/stiface"
	"google.golang.org/api/option"
)

const (
	objectName            = "develop/go.pdf.gpg"
	dstFilePath           = "go.pdf"
	gcsServiceAccountPath = "/Users/adamplansky/workspace/go/src/github.com/EndevelCZ/go-crypt/secrets/cterminal-developer.json"
	bucketName            = "cterminal-store-wallets"
	privFile              = "/Users/adamplansky/workspace/go/src/github.com/EndevelCZ/go-crypt/secrets/aliceprivtest.asc"
	passphraseFile        = "/Users/adamplansky/workspace/go/src/github.com/EndevelCZ/go-crypt/secrets/aliceprivtest-pass.txt"
)

// echo -n 'password' > /Users/adamplansky/workspace/go/src/github.com/EndevelCZ/go-crypt/secrets/aliceprivtest-pass.txt
func main() {
	passphraseByte, err := ioutil.ReadFile(passphraseFile)
	if err != nil {
		panic(fmt.Errorf("unable to open file: %s %s", passphraseFile, err))
	}
	// ------ open priv ------
	f, err := os.Open(privFile)
	if err != nil {
		panic(fmt.Errorf("unable to open file: %s %s", privFile, err))
	}
	defer f.Close()

	entityList, err := crypto.ReadArmorKeyring(f)
	if err != nil {
		panic(fmt.Errorf("unable to ReadArmorKeyring %s %#v", f, entityList))
	}

	entity := entityList[0]
	entity.PrivateKey.Decrypt(passphraseByte)
	for _, subkey := range entity.Subkeys {
		err := subkey.PrivateKey.Decrypt(passphraseByte)
		if err != nil {
			panic(fmt.Errorf("unable to decrypt private key: %s", err))
		}
	}
	// ------ open priv ------

	// ------ create client ------
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(gcsServiceAccountPath))
	if err != nil {
		panic(fmt.Sprintf("unable to create gcs client: %v", err))
	}
	gcs := gcp.NewClient(stiface.AdaptClient(client), bucketName, ctx)
	// ------ create client ------

	// ------ download ------
	src := &bytes.Buffer{}
	err = gcs.DownloadObject(src, objectName)
	if err != nil {
		panic(fmt.Sprintf("unable to download file: %v", err))
	}
	fmt.Println("Download OK")
	// ------ download ------
	var b []byte
	dst := bytes.NewBuffer(b)
	err = crypto.Decrypt(src, dst, entityList)
	if err != nil {
		panic(fmt.Errorf("unable to decrypt string err: %s", err))
	}
	// fmt.Println(dst.String())

	err = ioutil.WriteFile(dstFilePath, dst.Bytes(), 0644)
	if err != nil {
		panic(fmt.Errorf("unable to write to file: %s %s", dstFilePath, err))
	}

}
