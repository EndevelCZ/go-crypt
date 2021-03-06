package main

import (
	"bytes"
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
	"google.golang.org/api/option"
)

const (
	objectName            = "EncryptFn.txt"
	gcsServiceAccountPath = "/Users/adamplansky/workspace/go/src/github.com/EndevelCZ/go-crypt/secrets/cterminal-developer.json"
	bucketName            = "cterminal-store-wallets"
	// encryptionFilePath    = "/Users/adamplansky/Desktop/books/golang/1-The.Go.Programming.Language.pdf"
	encryptionFilePath = "/Users/adamplansky/Desktop/books/golang/1-LEARNING-GO.txt"
	// encryptionFilePath    = "/Users/adamplansky/Desktop/books/golang/1-The.Go.Programming.Language.pdf"
	// encryptionFilePath = "/Users/adamplansky/Movies/retrainyourmind.mp4"

	pubArmorKeyAdamPlansky = `-----BEGIN PGP PUBLIC KEY BLOCK-----

mQINBFh0nQUBEADIymzQp91+kaLkEA699NOSuw0GDLHpEx73G2e9EDh/GO3MAT0O
MBMYPiPkGi1wmcygBjZGW8H17RAGbxVBthlsaVRu43NgO+BdPn4K3lBDVUMeYQDp
OW5xOBpGzhUNH68hIyRB+Qex4GR6Zm97MmMXMAUZ2Mqv9QCOZqMyZV7Zzp3NA2Gx
Qr8AQC9mHhCBaZ7/Wq1esC682U467RaPIMuV7JiFwpT2MdBa1T57s+yxaFB8qw0p
kqOdp/2ANoqmaNGzVk+gybM0ter9pQUbaaFsRk5nmwMn1ZUV8G1NfGvWKr+tyjED
RCZI/KHgEwt5KZ8A8+S+f1dbuCDZrQI+9wFIZgRmoK/WOuWx5DSXLuf7ivGRzt0L
4tcinuRhxucNp5nZXa+FnDHmlnNqvxlPwu8TLKMp0DBSXXNXOIZKW6r6RUHCim9B
K44HlNeLn2ygid7J3+bHjy9wEuiPmvWOPOLxWC8hV+UMGXf+ex0IThrQ7+Cmb5U4
GHujPLyNiFOmTSIbkYVLVrXSxyyJ3uN6/8n0wEbvEZSAPC0j5AkKYyET7zOumADj
q02yOjHUvJcNzQmikhztKlbWoFIXsTjhSo4orA3MveZJ8m+YL5BjL26PM87vsvqj
7luWEnS7h/MYbah9pUCmKZT3RyVOdKJy56gt+s+5TlFKERwDjtcxHby9gQARAQAB
tCRBZGFtIFBsYW5za3kgPGFkYW1wbGFuc2t5QGdtYWlsLmNvbT6JAjkEEwEIACMF
Alh0nQUCGwMHCwkIBwMCAQYVCAIJCgsEFgIDAQIeAQIXgAAKCRCyzXHIVTpTShFo
EACzfcBAiFAe1eN+Hj7+TZlL20segIjYPOr/YmW1w99BguiCm08V5XSlNGK+tzke
FYC51O9v20QqFDdCugrnN8tEcR42ujri6k4MvO2876kw9cRI2vPeGcz5hB8bPcH+
KE6/f2ml0JDJ4rPKSYvMXyLIj2OUTEhQuN7JO8eMfoCF2Lpx7d/tGl8CkwvJfv+q
P9Bszn//Mh+nuW7APfO+qvI1KxG4eKuXE7yCo3W/1o/gYh9uIGtWpMO7CzOvMTa1
EIJB1EJRsRNk0EAnieMZU81LrAyo+ECUQw8rJjxRnBtaeE0wD51Ishh17uafKZMQ
KQZu57oT0envFS/56KQIYd6ejXvg1VUq3Et4CzbUVmQRCnOmenYT9fUnf3meUo/q
i9r5tI4Usw4R3mTTYN6m3go6ggSCwWeUZERI7CTBnZwCe0Z/CiIqr5Ehdzf7qCAN
MSiO+ZX0u1l3O7KqjZ8cBnUeaJtUGn2SAHprpXs4usQlETwfZHroIL9Pj7XC1rL7
eM3GtJ7MNF4NSdUVMwS0zrHDqxGeUPBfQarOsb5q9lZnvestU6Yg+n9NmUbW4xQD
FBxdM2Or9V/kMiwmXSgz4JP5+HIxN4oyUpYhkoLTWjlnp/1mFAkXL1kylG8wSRaV
X091DN5tXiQiATJIrxgZ4Fofp6d+Iwqzqo4Ly9RlN8HWCrkCDQRYdJ0FARAAmWNo
f9Ickfzw9fSXNdjE9tvSo3fB5MgShUeWQbqA4RjLn+QvBXncM4oWYVstD9bPzR+S
+/YX/2Djv3agkWzxzarfJ3Zn8FTQAtqJxQzBeHsHr7qL2zm+NmO32w8N/XSS3qsY
HHLPnAQUI9tmRdMTnukLJwXZGh/01FsdFny0676fnLMS3XpJBSjeowPR5G4SeDd0
dcyTiXIz/wX3/dDMxLXGjgDUp7qRKqk17uk4rYsSHzHiitmvwg9DtRO4lic1Cevr
uTx5nnusX9BshcqQ5PJsDA5Aq48f1hfLXB+zGXbFcvzg95C+lEmceTuhauGMc1TB
hhH7W/xEG+rNXAqrVz7xgdCvAZaMOEpo9F3wsCQCsb5fnoN68Xqcj1A2c7TnKtjs
fN4gBhjTqYd1CT7JPnW2sKk91s387GB01BrB5pmieVygP7hCnFtTL9g0UE71i/ns
u89JBgV7mtgLE/IWMp4LEAIhCOnFx8F+fy9PYsTXBzZS8PHtanEegGvjAq63u/Sc
10ShNfOfN8gtGfNVDcb2k9k+4rn3XmsZ3vkz37YFGkJkOz/9Zm4OQggCwVneXpac
opMZV+t29ynsBHkEyDCs8T7Zjwd5scYLChdRAO94NorN8LOjkehtiXMDnW7zSAtI
UrS51bI8HfJNMX54I0syuSL+yr5E3eA7TaD5cXsAEQEAAYkCHwQYAQgACQUCWHSd
BQIbDAAKCRCyzXHIVTpTSm0xEAC/rlulxwr45lTqz8H5XdoQG9fApoYMswbYIu1a
vESkamRZka6no1A+/Nb7pTu4k0qEpUtGu6Xy47OgGrG+jVpxCwnSm9Yewk0Ep2JW
kPW/SUK9+0+vrKPf2xjfrTHspIFxUH8VCX+kVrdA5xTuU/RFwcsyRvs4DBLKfDEm
YT7E9yrFYpQBUog7EmbyG3wc2wnAq1iuEeTzh4EjGmEh+kd3CQblyoR2a2f/mYXg
g+pEkGGS/+GxUc+4UgLh+s1S/VX0dnU3cpVw0iUimvelqMgFD/Qbu5CBSmvAAZXx
MTRzaNSoeUZQczU1XfDmDblxoKkI6YJ9hORnkbHEUWU0q+xRxBBfRitOgptXJPIv
uOL9Ml5Itul6O4ItwBuqeDn8oXMBMthoHsOQbKUuY7bO56669WGaQa/LhmZ/JDg/
opUXLGz83rpRjA0OWlI5CXfNQq0i8yT0OzTBOtIPf6eMS5O6bZXCrIsb7JCvV35k
gsUzxBYEwDYGYVbyeD74ch8q2OxKbP2ohAPAYznNia58qTu1qbaTu/JsYKypujZM
e5eul/dA2kNWrJc7IY8XMMW3nlwlUaLvEEHNsIFzNCUe8xwrG7hFEHdVxx68myDW
ns7iqH9oaSmsGLJ1sRX2gph9albI06fgCE8t2SFvfQbQqzHPzmQJHY9MW23zqeOe
8aQovQ==
=uPAR
-----END PGP PUBLIC KEY BLOCK-----`
)

var (
	pubKeys = []io.Reader{
		bytes.NewBufferString(pubArmorKeyAdamPlansky),

		// os.Open("/Users/adamplansky/workspace/go/src/github.com/EndevelCZ/go-crypt/secrets/karelonprem.asc"),
		// os.Open("/Users/adamplansky/workspace/go/src/github.com/EndevelCZ/go-crypt/secrets/alicepubtest.asc"),
	}
)

func main() {

	f, err := os.Open("/Users/adamplansky/workspace/go/src/github.com/EndevelCZ/go-crypt/secrets/aliceonprem.asc")
	if err != nil {
		panic(fmt.Sprintf("unable to open: %v", err))
	}
	defer f.Close()
	pubKeys = append(pubKeys, f)

	defer profile.Start().Stop()
	var mem runtime.MemStats
	log.Println("memory baseline...")

	runtime.ReadMemStats(&mem)

	// ------ create client ------
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(gcsServiceAccountPath))
	if err != nil {
		panic(fmt.Sprintf("unable to create gcs client: %v", err))
	}
	gcs := gcp.NewClient(stiface.AdaptClient(client), bucketName, ctx)

	r, err := os.Open(encryptionFilePath)
	if err != nil {
		panic(fmt.Errorf("unable to open file: %s ", err))
	}
	defer r.Close()

	err = gcs.UploadEncyptObject(r, objectName, crypto.EncryptFn, pubKeys...)
	if err != nil {
		panic(fmt.Errorf("unable to encrypt wc or upload: %s ", err))
	}

	fmt.Println("encrypted file OK!")
	fmt.Println("upload file OK!")
	log.Println("memory comparison...")

	runtime.ReadMemStats(&mem)
	log.Println("mem alloc: ", bytefmt.ByteSize(uint64(mem.Alloc)))
	log.Println("mem total alloc: ", bytefmt.ByteSize(uint64(mem.TotalAlloc)))
	log.Println("mem heap alloc: ", bytefmt.ByteSize(uint64(mem.HeapAlloc)))
	log.Println("mem heap sys: ", bytefmt.ByteSize(uint64(mem.HeapSys)))

}
