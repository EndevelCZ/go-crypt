package main

import (
	"bytes"
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"github.com/EndevelCZ/go-crypt/pkg/gcp"
	"github.com/googleapis/google-cloud-go-testing/storage/stiface"
	"google.golang.org/api/option"
)

const (
	filename              = "develop/fml.txt"
	gcsServiceAccountPath = "/Users/adamplansky/workspace/go/src/github.com/EndevelCZ/go-crypt/secrets/cterminal-developer.json"
	bucketName            = "cterminal-store-wallets"
)

func main() {
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
