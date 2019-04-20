package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/EndevelCZ/go-crypt/internal/interface/grpc"
)

func check(err error) {
	if err == nil {
		return
	}

	fmt.Printf("ERROR: %+v\n", err)
	os.Exit(1)
}
func main() {
	grpcClient, err := grpc.NewGcsClientGRPC("localhost:5000")
	check(err)
	// err = grpcClient.UploadFile(context.Background())
	f, err := os.Open("/Users/adamplansky/workspace/go/src/github.com/EndevelCZ/go-crypt/cmd/cterminal-store/upload/1.jpg")
	if err != nil {
		panic(fmt.Errorf("cannot open file %s", err))
	}
	defer f.Close()
	// f := strings.NewReader("hello world")
	err = grpcClient.Upload(context.Background(), f)
	if err != nil {
		logrus.Errorln(err)
	}
}
