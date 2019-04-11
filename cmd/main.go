package main

import (
	"fmt"
	"log"
	"os"

	"github.com/getsentry/raven-go"
)

func init() {
	raven.SetDSN("https://7a631e0a955249b8b05fcf3d04350ef2:e4f9c2e05b5e4448855283272d5e9312@sentry.io/1435605")
}

func main() {
	f, err := os.Open("filename.ext")
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.Panic(err)
	}
	defer f.Close()
	fmt.Println(f)
}
