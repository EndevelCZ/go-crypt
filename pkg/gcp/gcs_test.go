package gcp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/googleapis/google-cloud-go-testing/storage/stiface"
)

func TestUploadObject(t *testing.T) {
	bucketName := "my bucket"
	ctx := context.Background()
	client := newFakeClient()
	bkt := client.Bucket(bucketName)
	if err := bkt.Create(ctx, "my-project", nil); err != nil {
		t.Fatal(err)
	}
	gcs := NewClient(client, bucketName, ctx)
	type args struct {
		r          io.Reader
		objectName string
	}
	tests := []struct {
		name    string
		want    string
		args    args
		wantErr bool
	}{
		{
			name: "upload test",
			want: "hello",
			args: args{
				r:          strings.NewReader("hello"),
				objectName: "tester",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := gcs.UploadObject(tt.args.r, tt.args.objectName); (err != nil) != tt.wantErr {
				t.Errorf("UploadObject() error = %v, wantErr %v", err, tt.wantErr)
			}
			r, err := bkt.Object(tt.args.objectName).NewReader(ctx)
			if err != nil {
				t.Errorf("unable to read object")
			}
			sb := bytes.NewBufferString("")
			if _, err := io.Copy(sb, r); err != nil {
				t.Errorf("io copy error")
			}
			got := sb.String()
			if got != tt.want {
				t.Errorf("upload error got: %s, want: %s", got, tt.want)
			}
		})
	}
}

func TestDownloadObject(t *testing.T) {
	bucketName := "my bucket"
	ctx := context.Background()
	client := newFakeClient()
	bkt := client.Bucket(bucketName)
	if err := bkt.Create(ctx, "my-project", nil); err != nil {
		t.Fatal(err)
	}

	gcs := NewClient(client, bucketName, ctx)
	type args struct {
		w          io.Writer
		objectName string
	}
	tests := []struct {
		name    string
		want    string
		args    args
		wantErr bool
	}{
		{
			name: "download test",
			want: "hello",
			args: args{
				w:          &bytes.Buffer{},
				objectName: "tester",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ww := bkt.Object(tt.args.objectName).NewWriter(ctx)
			sb := bytes.NewBufferString(tt.want)
			if _, err := io.Copy(ww, sb); err != nil {
				t.Errorf("io copy error")
			}
			_ = ww.Close()

			w := &bytes.Buffer{}
			if err := gcs.DownloadObject(w, tt.args.objectName); (err != nil) != tt.wantErr {
				t.Errorf("DownloadObject() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := w.String(); got != tt.want {
				t.Errorf("DownloadObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestListObjects(t *testing.T) {
// 	type args struct {
// 		projectID          string
// 		bucket             string
// 		serviceAccountPath string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantW   string
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			w := &bytes.Buffer{}
// 			if err := ListObjects(w, tt.args.projectID, tt.args.bucket, tt.args.serviceAccountPath); (err != nil) != tt.wantErr {
// 				t.Errorf("ListObjects() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if gotW := w.String(); gotW != tt.wantW {
// 				t.Errorf("ListObjects() = %v, want %v", gotW, tt.wantW)
// 			}
// 		})
// 	}
// }

type fakeClient struct {
	stiface.Client
	buckets map[string]*fakeBucket
}

type fakeBucket struct {
	attrs   *storage.BucketAttrs
	objects map[string][]byte
}

func newFakeClient() stiface.Client {
	return &fakeClient{buckets: map[string]*fakeBucket{}}
}

func (c *fakeClient) Bucket(name string) stiface.BucketHandle {
	return fakeBucketHandle{c: c, name: name}
}

var _ stiface.BucketHandle = fakeBucketHandle{}

type fakeBucketHandle struct {
	stiface.BucketHandle
	c    *fakeClient
	name string
}

// type fakeObjectIterator struct {
// 	stiface.ObjectIterator
// 	b *fakeBucketHandle
// }

func (b fakeBucketHandle) Create(_ context.Context, _ string, attrs *storage.BucketAttrs) error {
	if _, ok := b.c.buckets[b.name]; ok {
		return fmt.Errorf("bucket %q already exists", b.name)
	}
	if attrs == nil {
		attrs = &storage.BucketAttrs{}
	}
	attrs.Name = b.name
	b.c.buckets[b.name] = &fakeBucket{attrs: attrs, objects: map[string][]byte{}}
	return nil
}

func (b fakeBucketHandle) Attrs(context.Context) (*storage.BucketAttrs, error) {
	bkt, ok := b.c.buckets[b.name]
	if !ok {
		return nil, fmt.Errorf("bucket %q does not exist", b.name)
	}
	return bkt.attrs, nil
}

func (b fakeBucketHandle) Object(name string) stiface.ObjectHandle {
	return fakeObjectHandle{c: b.c, bucketName: b.name, name: name}
}

// func (b fakeBucketHandle) Objects(context.Context, *storage.Query) stiface.ObjectIterator {
// 	return &{

// 	}
// }

type fakeObjectHandle struct {
	stiface.ObjectHandle
	c          *fakeClient
	bucketName string
	name       string
}

func (o fakeObjectHandle) NewReader(context.Context) (stiface.Reader, error) {
	bkt, ok := o.c.buckets[o.bucketName]
	if !ok {
		return nil, fmt.Errorf("bucket %q not found", o.bucketName)
	}
	contents, ok := bkt.objects[o.name]
	if !ok {
		return nil, fmt.Errorf("object %q not found in bucket %q", o.name, o.bucketName)
	}
	return fakeReader{r: bytes.NewReader(contents)}, nil
}

func (o fakeObjectHandle) Delete(context.Context) error {
	bkt, ok := o.c.buckets[o.bucketName]
	if !ok {
		return fmt.Errorf("bucket %q not found", o.bucketName)
	}
	delete(bkt.objects, o.name)
	return nil
}

type fakeReader struct {
	stiface.Reader
	r *bytes.Reader
}

func (r fakeReader) Read(buf []byte) (int, error) {
	return r.r.Read(buf)
}

func (r fakeReader) Close() error {
	return nil
}

func (o fakeObjectHandle) NewWriter(context.Context) stiface.Writer {
	return &fakeWriter{obj: o}
}

type fakeWriter struct {
	stiface.Writer
	obj fakeObjectHandle
	buf bytes.Buffer
}

func (w *fakeWriter) Write(data []byte) (int, error) {
	return w.buf.Write(data)
}

func (w *fakeWriter) Close() error {
	bkt, ok := w.obj.c.buckets[w.obj.bucketName]
	if !ok {
		return fmt.Errorf("bucket %q not found", w.obj.bucketName)
	}
	bkt.objects[w.obj.name] = w.buf.Bytes()
	return nil
}

// func Test_gcs_UploadEncyptObject(t *testing.T) {
// 	type fields struct {
// 		client     stiface.Client
// 		bucketName string
// 		ctx        context.Context
// 	}
// 	type args struct {
// 		r          io.Reader
// 		objectName string
// 		fn         encryptFn
// 		recip      []string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "hello",
// 			field: field{
// 				client:     newFakeClient(),
// 				bucketName: "cterminal-store-wallets",
// 				ctx:        context.Background(),
// 			},
// 			args: args{
// 				r:          bytes.NewBufferString("hello world"),
// 				objectName: "helloWorld.txt",
// 				fn:         func(w io.Writer, recip ...string) (io.WriteCloser, error) {
// 					bytes.New
// 				},
// 				recip:      nil,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			gcs := &gcs{
// 				client:     tt.fields.client,
// 				bucketName: tt.fields.bucketName,
// 				ctx:        tt.fields.ctx,
// 			}
// 			if err := gcs.UploadEncyptObject(tt.args.r, tt.args.objectName, tt.args.fn, tt.args.recip...); (err != nil) != tt.wantErr {
// 				t.Errorf("gcs.UploadEncyptObject() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
