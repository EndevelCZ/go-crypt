package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"cloud.google.com/go/storage"
	"github.com/EndevelCZ/go-crypt/pkg/crypto"
	"github.com/EndevelCZ/go-crypt/pkg/gcp"
	"github.com/pkg/profile"

	"code.cloudfoundry.org/bytefmt"
	"github.com/googleapis/google-cloud-go-testing/storage/stiface"
	"golang.org/x/crypto/openpgp"
	"google.golang.org/api/option"
)

const (
	objectName            = "learning-go-better.txt"
	gcsServiceAccountPath = "/Users/adamplansky/workspace/go/src/github.com/EndevelCZ/go-crypt/secrets/cterminal-developer.json"
	bucketName            = "cterminal-store-wallets"
	// encryptionFilePath    = "/Users/adamplansky/Desktop/books/golang/1-The.Go.Programming.Language.pdf"
	// encryptionFilePath = "/Users/adamplansky/Desktop/books/golang/1-LEARNING-GO.txt"
	// encryptionFilePath    = "/Users/adamplansky/Desktop/books/golang/1-The.Go.Programming.Language.pdf"
	encryptionFilePath = "/Users/adamplansky/Movies/retrainyourmind.mp4"
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
	defer profile.Start().Stop()
	var mem runtime.MemStats
	log.Println("memory baseline...")

	runtime.ReadMemStats(&mem)
	// ------ export public keys ------
	entities, err := crypto.ReadEntitesFromFiles(pubKeys...)
	if err != nil {
		panic(fmt.Errorf("unable to read pubkeys from file %s", err))
	}
	// ------ export public keys ------

	// var buf bytes.Buffer
	// b, err := ioutil.ReadFile(encryptionFilePath)
	// src := bytes.NewBuffer(b)
	// // src := bytes.NewBufferString(encryptionText)
	// err = crypto.Encrypt(entities, nil, src, &buf)

	// r, err := crypto.EncryptFile(entities, encryptionFilePath)
	// if err != nil {
	// 	panic(fmt.Errorf("unable to encrypt string: %s ", err))
	// }
	// fmt.Println("encrypted file OK")

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
	// s := &buf

	r, err := os.Open(encryptionFilePath)
	if err != nil {
		panic(fmt.Errorf("unable to open file: %s ", err))
	}
	defer r.Close()

	wcGcs := gcs.UploadObjectWriter(objectName)
	defer wcGcs.Close()

	// ------ upload ------
	wc, err := openpgp.Encrypt(wcGcs, entities, nil, &openpgp.FileHints{IsBinary: true}, nil)
	if err != nil {
		panic(fmt.Errorf("unable to encrypt wc: %s ", err))
	}
	if _, err := io.Copy(wc, r); err != nil {
		panic(fmt.Errorf("unable to copy encrypted bytes to gcs: %s ", err))
	}
	if err := wc.Close(); err != nil {
		panic(fmt.Errorf("unable to close gcs file: %s ", err))
	}
	fmt.Println("encrypted file OK!")
	fmt.Println("upload file OK!")
	log.Println("memory comparison...")

	runtime.ReadMemStats(&mem)
	log.Println("mem alloc: ", bytefmt.ByteSize(uint64(mem.Alloc)))
	log.Println("mem total alloc: ", bytefmt.ByteSize(uint64(mem.TotalAlloc)))
	log.Println("mem heap alloc: ", bytefmt.ByteSize(uint64(mem.HeapAlloc)))
	log.Println("mem heap sys: ", bytefmt.ByteSize(uint64(mem.HeapSys)))
	// b := bytefmt.ByteSize(uint64(mem.TotalAlloc)) // "1K"
	// log.Println("total alloc: ", b)
}
