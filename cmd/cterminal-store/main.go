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
	filename              = "develop/fml.txt"
	gcsServiceAccountPath = "/Users/adamplansky/workspace/go/src/github.com/EndevelCZ/go-crypt/secrets/cterminal-developer.json"
	bucketName            = "cterminal-store-wallets"
)

var (
	pubKeys = []string{
		"adamplansky.asc",
		"aliceonprem.asc",
		// "karelonprem.asc",
	}
)

func main() {
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

	// dst := bytes.NewBufferString("")

	dstFile := "encrypted1.gpg"
	dst, err := os.Create(dstFile)
	if err != nil {
		panic(err)
	}
	defer dst.Close()

	src := bytes.NewBufferString("this is encrypted string for multiple parties!")
	// err = Encrypt([]*openpgp.Entity{recipient}, nil, src, dst)
	err = crypto.Encrypt(entities, nil, src, dst)
	if err != nil {
		panic(fmt.Errorf("unable to encrypt string: %s ", err))
	}
	fmt.Println("encrypted string: ", dstFile)

	// recipient, err := ReadEntity(bytes.NewBufferString(pubArmorKey))
	return
	// ------ create client ------
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(gcsServiceAccountPath))
	if err != nil {
		panic(fmt.Sprintf("unable to create gcs client: %v", err))
	}
	gcs := gcp.NewClient(stiface.AdaptClient(client), bucketName, ctx)
	// ------ create client ------

	// ------ upload ------
	s := bytes.NewBufferString("hello")
	err = gcs.UploadObject(s, filename)
	if err != nil {
		panic(fmt.Sprintf("unable to upload file: %v", err))
	}
	// ------ upload ------

	// ------ download ------
	b := &bytes.Buffer{}
	err = gcs.DownloadObject(b, filename)
	if err != nil {
		panic(fmt.Sprintf("unable to download file: %v", err))
	}
	fmt.Println(b.String())
	// ------ download ------

}
