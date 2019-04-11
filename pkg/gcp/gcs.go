package gcp

import (
	"context"
	"fmt"
	"io"
	"log"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type gcs struct {
	bucketName, projectId string
	serviceAccountPath    string
}

func UploadGcsObject(r io.Reader, projectID, bucket, object, serviceAccountPath string) error {
	ctx := context.Background()
	// Creates a client.
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(serviceAccountPath))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, r); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	return nil
}

func DownloadGcsObject(w io.Writer, projectID, bucket, object, serviceAccountPath string) error {
	ctx := context.Background()
	// Creates a client.
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(serviceAccountPath))
	fmt.Println(serviceAccountPath)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return err
	}
	defer rc.Close()
	if _, err = io.Copy(w, rc); err != nil {
		return err
	}
	return nil
}
func ListGcsObjects(w io.Writer, projectID, bucket, serviceAccountPath string) error {
	ctx := context.Background()
	// Creates a client.
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(serviceAccountPath))
	if err != nil {
		return fmt.Errorf("Failed to create client: %v", err)
	}
	it := client.Bucket(bucket).Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintln(w, attrs.Name)
	}
	return nil
}
