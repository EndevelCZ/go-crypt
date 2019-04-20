package gcp

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/googleapis/google-cloud-go-testing/storage/stiface"
	"google.golang.org/api/iterator"
)

//https://medium.com/@tech_phil/how-to-stub-external-services-in-go-8885704e8c53
//https://itnext.io/how-to-stub-requests-to-remote-hosts-with-go-6c2c1db32bf2
//https://stackoverflow.com/questions/51799126/mock-external-dependencies-in-golang
//https://stackoverflow.com/questions/39814108/how-to-unit-test-google-cloud-storage-in-golang
//https://stackoverflow.com/questions/19167970/mock-functions-in-go
//https://aircto.com/blog/3-simple-and-powerful-tips-for-testing-golang-apps/
//https://stackoverflow.com/questions/54285068/how-to-mock-request-request-in-golang
//https://docs.aws.amazon.com/sdk-for-go/api/service/lambda/lambdaiface/
//https://stackoverflow.com/questions/39814108/how-to-unit-test-google-cloud-storage-in-golang
//https://github.com/aws/aws-sdk-go/tree/master/example/service/sqs/mockingClientsForTests
//!!!
//https://github.com/googleapis/google-cloud-go/issues/592#issuecomment-406099221
//https://godoc.org/cloud.google.com/go/httpreplay
//https://github.com/googleapis/google-cloud-go-testing
//https://github.com/googleapis/google-cloud-go/blob/516560fabc7a8e8b6f53a2e0a2b0d74d1b136f0f/storage/storage_test.go

// type gcs struct {
// 	bucketName, projectID string
// 	serviceAccountPath    string
// }

// NewGcs create new gcs object
// func NewGcs(bucketName, projectID, serviceAccountPath string) Storager {
// 	return &gcs{
// 		bucketName:         bucketName,
// 		projectID:          projectID,
// 		serviceAccountPath: serviceAccountPath,
// 	}
// }
type Gcs interface {
	DownloadObject(w io.Writer, objectName string) error
	UploadObject(r io.Reader, objectName string) error
	ListObjects(w io.Writer) error
	UploadObjectWriter(objectName string) io.WriteCloser
	UploadEncyptObject(r io.Reader, objectName string, fn encryptFn, recip ...io.Reader) error
}

type gcs struct {
	client     stiface.Client
	bucketName string
	ctx        context.Context
}

func NewClient(client stiface.Client, bucketName string, ctx context.Context) Gcs {
	return &gcs{
		client:     client,
		bucketName: bucketName,
		ctx:        ctx,
	}
}

func (gcs *gcs) UploadObjectWriter(objectName string) io.WriteCloser {
	wc := gcs.client.Bucket(gcs.bucketName).Object(objectName).NewWriter(gcs.ctx)
	return wc
}

type encryptFn func(w io.Writer, recip ...io.Reader) (io.WriteCloser, error)

func (gcs *gcs) UploadEncyptObject(r io.Reader, objectName string, fn encryptFn, recip ...io.Reader) error {
	w := gcs.client.Bucket(gcs.bucketName).Object(objectName).NewWriter(gcs.ctx)
	defer w.Close()
	wc, err := fn(w, recip...)
	if err != nil {
		return err
	}
	defer wc.Close()
	if _, err := io.Copy(wc, r); err != nil {
		return fmt.Errorf("unable to copy encrypted bytes to gcs: %s ", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("unable to close gcs file: %s ", err)
	}
	logrus.Infof("UploadEncyptObject: %s %s\n", objectName, time.Now())
	return nil
}

// UploadObject upload io.Reader to gcs with objectName
func (gcs *gcs) UploadObject(r io.Reader, objectName string) error {
	// ctx := context.Background()
	// // Creates a client.
	// client, err := storage.NewClient(ctx, option.WithCredentialsFile(gcs.serviceAccountPath))
	// if err != nil {
	// 	// log.Fatalf("Failed to create client: %v", err)
	// 	return err
	// }
	wc := gcs.client.Bucket(gcs.bucketName).Object(objectName).NewWriter(gcs.ctx)
	if _, err := io.Copy(wc, r); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	return nil
}

// DownloadObject download objectName to io.Writer
func (gcs *gcs) DownloadObject(w io.Writer, objectName string) error {
	// ctx := context.Background()
	// // Creates a client.
	// client, err := storage.NewClient(ctx, option.WithCredentialsFile(gcs.serviceAccountPath))
	// if err != nil {
	// 	// log.Fatalf("Failed to create client: %v", err)
	// 	return err
	// }
	rc, err := gcs.client.Bucket(gcs.bucketName).Object(objectName).NewReader(gcs.ctx)
	if err != nil {
		return err
	}
	defer rc.Close()
	if _, err = io.Copy(w, rc); err != nil {
		return err
	}
	return nil
}

// ListObjects show all objects in gcs
// write to io.Writer
func (gcs *gcs) ListObjects(w io.Writer) error {
	// ctx := context.Background()
	// // Creates a client.
	// client, err := storage.NewClient(ctx, option.WithCredentialsFile(gcs.serviceAccountPath))
	// if err != nil {
	// 	// return fmt.Errorf("Failed to create client: %v", err)
	// 	return err
	// }
	it := gcs.client.Bucket(gcs.bucketName).Objects(gcs.ctx, nil)
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
