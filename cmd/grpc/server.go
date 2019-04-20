package main

import (
	"fmt"
	"os"

	"github.com/EndevelCZ/go-crypt/internal/interface/grpc"
)

func check(err error) {
	if err == nil {
		return
	}

	fmt.Printf("ERROR: %+v\n", err)
	os.Exit(1)
}

const (
	gcsServiceAccountPath = "/Users/adamplansky/workspace/go/src/github.com/EndevelCZ/go-crypt/secrets/cterminal-developer.json"
	bucketName            = "grpc-cterminal"
	port                  = 5000
)

func main() {
	// port int, bucketName, gcsServiceAccountPath string
	grpcServer, err := grpc.NewGcsServerGRPC(port, bucketName, gcsServiceAccountPath)
	check(err)
	// server := &grpcServer
	err = grpcServer.Listen()
	check(err)
}
