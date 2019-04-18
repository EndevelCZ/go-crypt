package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"cloud.google.com/go/storage"
	"code.cloudfoundry.org/bytefmt"
	"github.com/EndevelCZ/go-crypt/pkg/crypto"
	"github.com/EndevelCZ/go-crypt/pkg/gcp"
	"github.com/pkg/profile"

	"github.com/googleapis/google-cloud-go-testing/storage/stiface"
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
	b, err := ioutil.ReadFile(encryptionFilePath)
	src := bytes.NewBuffer(b)
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
	fmt.Println("encrypted file OK!")
	fmt.Println("upload file OK!")
	log.Println("memory comparison...")

	runtime.ReadMemStats(&mem)
	log.Println("mem alloc: ", bytefmt.ByteSize(uint64(mem.Alloc)))
	log.Println("mem total alloc: ", bytefmt.ByteSize(uint64(mem.TotalAlloc)))
	log.Println("mem heap alloc: ", bytefmt.ByteSize(uint64(mem.HeapAlloc)))
	log.Println("mem heap sys: ", bytefmt.ByteSize(uint64(mem.HeapSys)))
}
