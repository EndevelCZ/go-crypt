package crypto

import (
	"bytes"
	"encoding/base64"
	"io"
	"testing"

	"golang.org/x/crypto/openpgp"
)

var (
	// plaintextFile      = "hello.txt"
	plaintextString = "hello"
	// ciphertextFile     = "hello.txt.gpg"
	ciphertextString = ""
	pubFilename      = "alice.asc"
	// privFilename           = "alice-priv.asc"
	privArmoryFilename     = "alice-priv-armor.asc"
	passphraseString       = "password"
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
	pubArmorKeyAlice = `-----BEGIN PGP PUBLIC KEY BLOCK-----

mQGNBFy03CIBDADjRJyWg8DhdAO4Lsig8CYGR1b9jas1BDoZ8cVazy8xtmEJ2vbe
h1JwAPduTQWfeuiJAEQ9MpaIiaa4p2JxRdV1EhnRfBu4w1zq4huCko3uDLGcVHkW
nZ3DFXiOeVPz+btDfY8YG1txQe51aFTy1FfF2ddLoXBBy4u2z2Zf823hVX9fX0Du
ksAWu1ZHfi5M+7R9spPYRXLkWHpS+oobL7Jl2k2seENmT6q/x0gmtUge8oSn8+Nx
ByZ/h8OYlzLFaA6zKjQ15MPN+PYo5SiUSjhncfXhMpDoEhW1RmZWNOSBImFu+JOK
eJb1h+ViwCnkP+GJz4P815rKuF97rfVtkBJwJaKS3PvILwRtj1IOnPUYjdc3lljv
vCKsrIC03JG5LUd5Z5jOqfZbkZ/cxZlFXFBecAqTDBmhoUucF/3MV7o3pJODyyqq
ORcLZk4T68T+BkjIKfO1SdAb9Vhl0OdD0wI7jJt52Chai7hy7wpm1BeHnUssX9vZ
tfF0kE7WGb3N3MMAEQEAAbQxa2FyZWwgb25wcmVtIChrYXJlbCBvbiBwcmVtKSA8
a2FyZWxAY3Rlcm1pbmFsLmN6PokBzgQTAQgAOBYhBM84rdBWVeeKPErLEQy9HTOR
9aNMBQJctNwiAhsDBQsJCAcCBhUICQoLAgQWAgMBAh4BAheAAAoJEAy9HTOR9aNM
ZRQL/RWg8MMygvfTGqZltbnONAoMjrYhwAvGbClgpm2nIfiZ69bHDVigB1wiJCDA
88f2WMcmzv1Gl9NMdQ8Tydt5wslppa7NthDrLN+gXo8ByqJw+IP5SENPHvflPEJ1
Zrd+UXYh60PVkeh2CLZ8KMx09RjAxqArkglqbJjE1BjfI+ozlmgYgVuy7S7EeL7J
FikxiSKPkogBfjXAMHgcO4nF+MtFozJDWVoJ/w++wT7SZfPoj0PO4GzIyR6SHjYy
6ZuicjRYdwt27H9/VoWLdZ94KuPr4+Sl4QQcOcw9uAQDTGuY6tVY1DlmJjmwGwuq
rRKGqLnhFEM+XqixkYwb0AlPeRtZDGIs066M3LgXs6SuNzjTU0bNDrXpPSUA/q6t
E2OCKX6ImA+FkfIykfQHAErjPaHzaBqxrzrQDxpGFl3sfFo2+Wv7VcoU/eFXB0Vt
Nm2sPB71Dc/AFecl76LRjLXu5JdfAnx0cCVHpOnc4+LPYjdlOjNhjmOQ5n7W+EyW
ICvFnLkBjQRctNwiAQwAwcIkpK+7AiJyHhbgCmZDE05kFM4gtFm/BoxJIa8Trnkk
1R7Xj4Z1mMiFBU6gIhmM3ccQzlgYonASGkpfA+oafRAFsOG/IR1Lm645/lRhXqD8
v+ipDY+HlblQDaUD+EvqLqKa4fXIlxwzj93SiVEk8lR4lIXqCUOyuAd1aWg7YWIO
OTE20HD4mcL/pZcfEhXA8F0mUB/MpAXn1tAQFAAI+44sQWs4oXNW/HmCKP/glaeo
+napngp+hTpXElWYb2xce1dTknNeq5noxC6MKNQWGfpHEoLSkboXhN410TcPeQMt
vJ4ODWU2TcenwFcUJCLQULKPtCJlh6o5HO1MOAV2ldvI0WdzwFEXnsVyewOlGyhe
PUQKJwOstKeB+k6Pd+sTGhAElNQnDczsIRKUQEr27tDNk6pkX8KS0Kvxx6L98u9h
B9AbUZ8kS/bHvePFDsJCSeitqZokPu2U9GgiwBqRET3l9ueToE0RMkr0qan2AQ2/
y6banw3BqwjErQo2HlT9ABEBAAGJAbYEGAEIACAWIQTPOK3QVlXnijxKyxEMvR0z
kfWjTAUCXLTcIgIbDAAKCRAMvR0zkfWjTDhIC/wJkW0VtIMPla558PC9l7XbWy5b
Z9CUstFKfmOCjkp6SW6YRtxewfTel/y5ohpdU3Vz2eQt/YVTmx2+CVb7wRLzxeMC
DA+BVi6ouzsNZR7XZL17eNyVzRXQTjDfEAoV/uPTBuVR530A2IcIO0jE5zTNAQ17
Mfmyxga+tF2A2SNC08d6nhjgFGiSIGtI1pnISPN77Ft+S3L60lL3EFC/jQ5131Pb
0U2Fot2+MfxhGv68swOKbS2uRy5h5Mip7ULACb9g0J2x1jj2HFP5rLETetLExFEA
qTN6CDjGdvNXYWnX5TPWDuLxnMoOmzvuOrz0CwZS8lSOAmgTTQtpuRdZqpTHuuUb
3rjNGIJIY4mmxAVCK6NSVJ2WSCGFH1HKB1aAmsKQFkepMe/1Tl9yvQfMVx8cHqND
CXYcZkFHnmKPkLE+NKj/WdekPh+wLquxgmvhzJC9Vlug5aIQT9OoA5uK0g6Pr0/g
ss4vBBVbL0DMrIXmbg8iZ9cOV9ez/qyLhNKaWQo=
=ZwzX
-----END PGP PUBLIC KEY BLOCK-----`
	pubArmorKey = `-----BEGIN PGP PUBLIC KEY BLOCK-----

mI0EXKG1jgEEAMoYJ+ID5dww0y8u5eOW28ZuoypreviYVHX7Qq1E6L/42birfqNt
+A89DcXfAuW2p3NEGqhndbp7M0+g6KBAmbKyXjdq22hH16vdp6hfGd6jxlD7wZM1
e78L9rB37P6XQklk1dmLII1/m1OW13RARgLOinqzdkJptgKJJQ9rnUWtABEBAAG0
L2FsaWNlQGN0ZXJtaW5hbC5jeiAoYWxpZWMpIDxhbGljZUBjdGVybWluYWwuY3o+
iM4EEwEIADgWIQTb+u5eUFdCX+wYZa8N4Ioxvr1okgUCXKG1jgIbAwULCQgHAgYV
CAkKCwIEFgIDAQIeAQIXgAAKCRAN4Ioxvr1okljCA/9v45oRNS1oYon/pvkmYTh5
ucCrwztFK9GDFBgl6yibdyVCPUWix3dZHvtXyDX8CEA0Tekv5eush38uC/09yhAE
Abcgcw47hH8DVsNBBOFJsvL3dPLxPE7n9wdMTlEQys2pewGNH6Ns4WQ3DZ42+r32
rngxmqmjDusl4fWlckC+XriNBFyhtY4BBADbsTBtTNT/0DOH3d3QSEKUpVWCRxT4
ZkiDeA84ag6ORrlTafQYcTW8K5BUKfXvvBQCSQlnVfhYtvXE81gdNDww3vrFwFkf
MwYa1f+VNI7tRA3390/VWq/505Yf4guWD2/KcVltMCbiXhqoXi3Ezq1pyAiRNYps
NxWnmWvYz75qlQARAQABiLYEGAEIACAWIQTb+u5eUFdCX+wYZa8N4Ioxvr1okgUC
XKG1jgIbDAAKCRAN4Ioxvr1okskjA/9Tzc/EusvO8zvVbxFAvVI5W/GTG4k/NQJv
8bVbY61sTfIomXglnrFluMWDmEVC6rgQZwN1gKawh+DUxFIm8pM++M0YI+iO5s1T
AOzZVxc0y1VX13jN9YG9GDBPJZlSE2GiACPiLZtLft1EIqZkmuR9NjvGqFr4CH/9
UcEGq+iQvg==
=4nyp
-----END PGP PUBLIC KEY BLOCK-----`
	privArmorKey = `-----BEGIN PGP PRIVATE KEY BLOCK-----

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
=swUV
-----END PGP PRIVATE KEY BLOCK-----`
)

// cat hello.txt |  tr -d \\n | base64
func TestGpgEncryptStringSymmetric(t *testing.T) {
	passphrase := []byte(passphraseString)
	fp := bytes.NewBufferString(plaintextString)
	fc := bytes.NewBufferString("")
	err := GpgEncryptSymmetric(fp, fc, passphrase)
	if err != nil {
		t.Errorf("unable to encrypt symmetric %s \n", err)
	}
	if len(fc.String()) < 1 {
		t.Errorf("symmetric encrypt failed %s\n", fc.String())
	}
}
func TestGpgDecryptSymmetricString(t *testing.T) {
	passphrase := []byte(passphraseString)
	symmetricGpgTextBase64 := "wx4EBwMIEq4Yod9IW3pg4NX4O/wVtcuj7Ndr8J2U9gXS4AHkRv2mpwVPrtcyj1bI9K/lsOGNBuBC4Gvhdrjg3eJd6LKW4NfiGro3tuCK4A/gCOQmUm9bl+MqpaQ7TI5ds51r4u2D8/Th65YA"
	sDec, _ := base64.StdEncoding.DecodeString(symmetricGpgTextBase64)
	fp := bytes.NewBufferString("")
	fc := bytes.NewBuffer(sDec)

	err := GpgDecryptSymmetric(fc, fp, passphrase)
	if err != nil {
		t.Errorf("unable to decrypt file %s\n", err)
	}
	if len(fp.String()) < 1 {
		t.Errorf("symmetric decrypt failed %s\n", fp.String())
	}

	want := plaintextString
	got := fp.String()
	if fp.String() != plaintextString {
		t.Errorf("unable to decrypt string want: %s\n got: %s", want, got)
	}
}
func TestGpgEncryptDecryptSymmetric(t *testing.T) {
	passphrase := []byte(passphraseString)
	fp := bytes.NewBufferString(plaintextString)
	fc := bytes.NewBufferString("")
	err := GpgEncryptSymmetric(fp, fc, passphrase)
	if err != nil {
		t.Errorf("unable to encrypt symmetric %s \n", err)
	}
	fd := bytes.NewBufferString("")
	if len(fc.String()) < 1 {
		t.Errorf("symmetric encrypt failed %s\n", fc.String())
	}

	// echo -n 'wx4EBwMIEq4Yod9IW3pg4NX4O/wVtcuj7Ndr8J2U9gXS4AHkRv2mpwVPrtcyj1bI9K/lsOGNBuBC4Gvhdrjg3eJd6LKW4NfiGro3tuCK4A/gCOQmUm9bl+MqpaQ7TI5ds51r4u2D8/Th65YA' | base64 -D > file_decrypted.gpg
	// x := fc.String()
	// sDec := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s", x)))
	// fmt.Printf("sDec: [%s]\n", sDec)

	err = GpgDecryptSymmetric(fc, fd, passphrase)
	if err != nil {
		t.Errorf("unable to decrypt file %s\n", err)
	}
	if len(fd.String()) < 1 {
		t.Errorf("symmetric decrypt failed %s\n", fd.String())
	}

	want := plaintextString
	got := fd.String()
	if fd.String() != plaintextString {
		t.Errorf("unable to decrypt string want: %s\n got: %s", want, got)
	}
}

func TestGpgEncryptStringSymetric(t *testing.T) {
	passphrase := []byte(passphraseString)

	pBuf := bytes.NewBufferString(plaintextString)
	cBuf := bytes.NewBufferString(ciphertextString)
	err := GpgEncryptSymmetric(pBuf, cBuf, passphrase)
	if err != nil {
		t.Errorf("unable to encrypt string %s, err: %s\n", plaintextString, err)
	}
	// fmt.Println(cBuf.String())
	// sEnc := base64.StdEncoding.EncodeToString(cBuf.Bytes())
	// t.Errorf("error: %s", sEnc)
	if cBuf.String() == ciphertextString {
		t.Errorf("unable to encrypt string %s, err: %s\n", plaintextString, err)
	}
}

func TestGpgDecryptStringSymetric(t *testing.T) {
	ciphertextString := "wx4EBwMI9lOpQbkJG39gXG37o6MkRH4pBX7XgUQL52vS4AHkmiXx67nyXSbSzyFU3jLx/uGQxuCy4DvhgNngRuJkw7a+4MvgxuCw5OqeRURBsM2jAB+3/g6AZMri2WER0eGrXwA="
	plaintextString := "1"
	passphrase := []byte(passphraseString)
	x, _ := base64.StdEncoding.DecodeString(ciphertextString)
	cBuf := bytes.NewBuffer(x)
	pBuf := bytes.NewBufferString("")
	err := GpgDecryptSymmetric(cBuf, pBuf, passphrase)
	if err != nil {
		t.Errorf("unable to decrypt string %s, err: %s\n", ciphertextString, err)
	}
	if pBuf.String() != plaintextString {
		t.Errorf("unable to decrypt string got: [%s], expecting: [%s]\n", plaintextString, pBuf.String())
	}
}

func TestReadEntity(t *testing.T) {
	b := bytes.NewBufferString(pubArmorKeyAlice)
	_, err := ReadEntity(b)
	if err != nil {
		t.Errorf("read entity error: %s\n", err)
	}
}

// func TestReadKeyring(t *testing.T) {
// 	fPriv, err := os.Open(privFilename)
// 	if err != nil {
// 		t.Errorf("unable to open file: %s %s", privFilename, err)
// 	}
// 	defer fPriv.Close()
// 	entityList, err := ReadKeyring(fPriv)
// 	if err != nil {
// 		t.Errorf("unable to ReadKeyring %s %#v", privFilename, entityList)
// 	}
// 	entity := entityList[0].Identities["alice@cterminal.cz (aliec) <alice@cterminal.cz>"]
// 	if entity.UserId.Comment != "aliec" {
// 		t.Errorf("entity doesn't exist %#v \n", entity)
// 	}
// }

func TestReadKeyringFromString(t *testing.T) {
	fPriv := bytes.NewBufferString(privArmorKey)
	entityList, err := ReadArmorKeyring(fPriv)
	if err != nil {
		t.Errorf("unable to ReadArmorKeyring %s %#v", fPriv, entityList)
	}
	entity := entityList[0].Identities["alice@cterminal.cz (aliec) <alice@cterminal.cz>"]
	if entity == nil || entity.UserId.Comment != "aliec" {
		t.Errorf("entity doesn't exist %#v \n", entity)
	}
}

func TestEncryptFromString(t *testing.T) {
	plaintexString := "hello"
	recipient, err := ReadEntity(bytes.NewBufferString(pubArmorKey))
	if err != nil {
		t.Errorf("read entity error: %s\n", err)
	}
	dst := bytes.NewBufferString("")
	src := bytes.NewBufferString(plaintexString)
	err = Encrypt([]*openpgp.Entity{recipient}, nil, src, dst)
	if err != nil {
		t.Errorf("unable to encrypt string: %s ", err)
	}
	decryptedString := base64.StdEncoding.EncodeToString(dst.Bytes())
	if len(decryptedString) < 1 {
		t.Errorf("unable to encrypt string %s got: %s\n ", plaintexString, decryptedString)
	}
}
func TestDecryptString(t *testing.T) {
	encryptedString := "wYwDJvpFsMUbEcoBBACW0VwBMZmLbKkMAPwyaoXVoqFFoF/2w+9+bxXfu2u9Md6Y2Axw9OUIYsClvnMJ1NyNsju0rb4O09YGIvhr2l3vhwzzNf5TU7rpd+RyH5RVaP4BeX/2sSmGsod/hb9qob2XNwdGv4Ajox1W4PkG+IVNsNtajU4kaZVn7rxiBLttjNLgAeTLaLRRBRrPgGqc04zi0/ZB4ZHK4K3gy+GJGOAi4oerb5LgMOJ3sLeD4Hfg8+Dl5DKMFLPoxv5sBBRHdAjI1XziVNoCWuGFpQA="
	cipherString, _ := base64.StdEncoding.DecodeString(encryptedString)
	pass := "2rYiibBZ33f9BMxNE6TEX$H.X4#aPyc2g*(LtRAauiGsn}PJT4{9(J)Xsbe@4Jbr"
	passphraseByte := []byte(pass)
	fPriv := bytes.NewBufferString(privArmorKey)
	entityList, err := ReadArmorKeyring(fPriv)
	if err != nil {
		t.Errorf("unable to ReadArmorKeyring %s %#v", fPriv, entityList)
	}

	entity := entityList[0]
	entity.PrivateKey.Decrypt(passphraseByte)
	for _, subkey := range entity.Subkeys {
		err := subkey.PrivateKey.Decrypt(passphraseByte)
		if err != nil {
			t.Errorf("unable to decrypt private key: %s", err)
		}
	}

	src := bytes.NewBuffer(cipherString)
	dst := bytes.NewBufferString("")
	err = Decrypt(src, dst, entityList)
	if err != nil {
		t.Errorf("unable to decrypt string err: %s\n", err)
	}
	got := dst.String()
	want := plaintextString
	if got != want {
		t.Errorf("asymetric encryption strinfg failed got: %s, want: %s\n", got, want)
	}
}

func TestReadEntities(t *testing.T) {
	type args struct {
		r []io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []*openpgp.Entity
		wantErr bool
	}{
		{
			name: "multiple entities",
			args: args{
				r: []io.Reader{
					bytes.NewBufferString(pubArmorKeyAdamPlansky),
					bytes.NewBufferString(pubArmorKeyAlice),
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadEntities(tt.args.r...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadEntities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, entity := range got {
				if len(entity.Identities) < 1 {
					t.Errorf("Unable to read pub key to entities")
				}
			}
		})
	}
}
