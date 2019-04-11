package gcp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"
)

func serviceAccount() (string, error) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		return "", fmt.Errorf("chdir error")
	}
	return fmt.Sprintf("%s/%s/%s", dir, "secrets", "adamplansky-gcs-go-crypt.json"), nil
}

var (
	projectID  = "adamplansky-181911"
	bucketName = "go-crypt"
)

func TestUploadFilePath(t *testing.T) {
	filePath := "file.txt"
	object := "file_test.txt"
	serviceAccountPath, err := serviceAccount()
	if err != nil {
		t.Errorf("err: %s", err)
	}
	f, err := os.Open(filePath)
	if err != nil {
		t.Errorf("unable to open file %s, %v", filePath, err)
	}
	defer f.Close()
	err = UploadGcsObject(f, projectID, bucketName, object, serviceAccountPath)
	if err != nil {
		t.Errorf("unable to upload file %s to gcs %s", filePath, object)
	}
}

func TestUploadString(t *testing.T) {
	object := "file_string_test.txt"
	serviceAccountPath, err := serviceAccount()
	if err != nil {
		t.Errorf("err: %s", err)
	}
	a := "heloo from string"
	c := bytes.NewBufferString(a)
	err = UploadGcsObject(c, projectID, bucketName, object, serviceAccountPath)
	if err != nil {
		t.Errorf("unable to upload string %s to gcs %s", a, object)
	}
}
func TestDownloadString(t *testing.T) {
	s := ""
	ss := "Hello world. Encryption and Decryption testing.\n"
	object := "file.txt"
	buf := bytes.NewBufferString(s)
	serviceAccountPath, err := serviceAccount()
	if err != nil {
		t.Errorf("err: %s", err)
	}
	err = DownloadGcsObject(buf, projectID, bucketName, object, serviceAccountPath)
	if err != nil {
		t.Errorf("gcs download error: %s", err)
	}
	rs := buf.String()
	if strings.Compare(rs, ss) != 0 {
		t.Errorf("strings are not equal")
	}

}
func TestDownloadObject(t *testing.T) {
	filePath := "file_test.txt"
	object := "file_test.txt"
	serviceAccountPath, err := serviceAccount()
	if err != nil {
		t.Errorf("err: %s", err)
	}
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Errorf("unable to open file %s, %v", filePath, err)
	}

	err = DownloadGcsObject(f, projectID, bucketName, object, serviceAccountPath)
	if err != nil {
		t.Errorf("gcs download error: %s", err)
	}
	//err = f.Sync()
	if err := f.Close(); err != nil {
		t.Errorf("close error: %s", err)
	}
	dat, err := ioutil.ReadFile("file.txt")
	if err != nil {
		t.Errorf("unable to open file %s to gcs %s", filePath, object)
	}

	d, _ := ioutil.ReadFile(filePath)
	for i := 0; i < len(dat); i++ {
		if dat[i] != d[i] {
			t.Errorf("gcs file is not same as local file")
		}
	}
}
func TestListObjects(t *testing.T) {
	serviceAccountPath, err := serviceAccount()
	if err != nil {
		t.Errorf("err: %s", err)
	}
	var buf bytes.Buffer
	err = ListGcsObjects(&buf, projectID, bucketName, serviceAccountPath)
	fmt.Println("-------")
	s := buf.String()
	fmt.Println(s)
	fmt.Println("-------")
	if err != nil {
		t.Errorf("gcs list error: %s", err)
	}
}
