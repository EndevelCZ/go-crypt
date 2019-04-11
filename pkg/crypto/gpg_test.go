package crypto

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"golang.org/x/crypto/openpgp"
)

var (
	plaintextFile      = "file.txt"
	ciphertextFile     = "file_enc.txt.gpg"
	plaintextString    = "1"
	ciphertextString   = ""
	pubFilename        = "alice.asc"
	privFilename       = "alice-priv.asc"
	privArmoryFilename = "alice-priv-armor.asc"
	pubArmorKey        = `-----BEGIN PGP PUBLIC KEY BLOCK-----

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

func TestGpgEncryptFileSymmetric(t *testing.T) {
	passphrase := []byte("password")

	fp, err := os.Open(plaintextFile)
	if err != nil {
		t.Errorf("unable to create file %s %s\n", plaintextFile, err)
	}
	defer fp.Close()
	fc, err := os.Create(ciphertextFile)
	if err != nil {
		t.Errorf("unable to create file %s %s\n", ciphertextFile, err)
	}
	defer fc.Close()

	err = GpgEncryptSymmetric(fp, fc, passphrase)
	if err != nil {
		t.Errorf("unable to encrypt file %s %s\n", plaintextFile, err)
	}
	if _, err := os.Stat(ciphertextFile); os.IsNotExist(err) {
		t.Errorf("unable to encrypt file %s\n", plaintextFile)
	}
}

func TestGpgDecryptFileSymmetric(t *testing.T) {
	passphrase := []byte("password")

	fc, err := os.Open(ciphertextFile)
	if err != nil {
		t.Errorf("unable to open file %s %s\n", ciphertextFile, err)
	}

	plaintextFile := "file_decrypted.txt"
	fp, err := os.Create(plaintextFile)
	if err != nil {
		t.Errorf("unable to create file %s %s\n", plaintextFile, err)
	}
	err = GpgDecryptSymmetric(fc, fp, passphrase)
	if err != nil {
		t.Errorf("unable to decrypt file %s\n", err)
	}
	if _, err := os.Stat(plaintextFile); os.IsNotExist(err) {
		t.Errorf("unable to decrypt file %s\n", plaintextFile)
	}
}

func TestGpgDecryptFileSymmetric2(t *testing.T) {
	passphrase := []byte("password")

	fp, err := os.Open(plaintextFile)
	if err != nil {
		t.Errorf("unable to create file %s %s\n", plaintextFile, err)
	}
	defer fp.Close()
	fc, err := os.Create(ciphertextFile)
	if err != nil {
		t.Errorf("unable to create file %s %s\n", ciphertextFile, err)
	}
	defer fc.Close()

	err = GpgEncryptSymmetric(fp, fc, passphrase)
	_, err = fc.Seek(0, 0)

	plaintextFile = "file_decrypted.txt"
	fp, err = os.Create(plaintextFile)
	if err != nil {
		t.Errorf("unable to create file %s %s\n", plaintextFile, err)
	}
	err = GpgDecryptSymmetric(fc, fp, passphrase)
	if err != nil {
		t.Errorf("unable to decrypt file %s\n", err)
	}
	if _, err := os.Stat(plaintextFile); os.IsNotExist(err) {
		t.Errorf("unable to decrypt file %s\n", plaintextFile)
	}
}

func TestGpgEncryptStringSymetric(t *testing.T) {
	passphrase := []byte("password")

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
	passphrase := []byte("password")
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
	ascFilename := "alice.asc"
	f, err := os.Open(ascFilename)
	if err != nil {
		t.Errorf("unable to open file: %s %s", ascFilename, err)
	}
	defer f.Close()
	_, err = ReadEntity(f)
	if err != nil {
		t.Errorf("read entity error: %s\n", err)
	}
}

func TestReadKeyring(t *testing.T) {
	fPriv, err := os.Open(privFilename)
	if err != nil {
		t.Errorf("unable to open file: %s %s", privFilename, err)
	}
	defer fPriv.Close()
	entityList, err := ReadKeyring(fPriv)
	if err != nil {
		t.Errorf("unable to ReadKeyring %s %#v", privFilename, entityList)
	}
	entity := entityList[0].Identities["alice@cterminal.cz (aliec) <alice@cterminal.cz>"]
	if entity.UserId.Comment != "aliec" {
		t.Errorf("entity doesn't exist %#v \n", entity)
	}
}

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
	encryptedString := "wYwDJvpFsMUbEcoBBADGVjJ9O8s5HSDwseG8tUxzgy2BtKuUvuVZHZeeZPQMqYFl9IMiKTrLv0n3jScOZs6rxUY/f8w4tXBffSje/66hIHHPJlN2pmxBDmSUl0y9DpbPDoDVjPIf9VbwbJz4lzx0vam8kn66TLC5rRsjWJvh0S7OKXl3okhKFBARJfDBgdLgAeRsGAI8XyfWlIsNoaZS6uAV4a+u4Prg9OGvs+BK4k4glU7gyuKd1X/84Ffhivrg7+SdsiPqEaeBhwmZWkDB/cre4owbkO7hmGkA"
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
	d := base64.StdEncoding.EncodeToString([]byte(dst.String()))
	if encryptedString != d {
		t.Errorf("unable to encrypt string: \n%s\n expecting:\n%s\n ", d, encryptedString)
	}
}

func TestEncrypt(t *testing.T) {
	fileToEnc := "file.txt"
	// f, err := os.Open(pubFilename)
	// if err != nil {
	// 	t.Errorf("unable to open file: %s %s", pubFilename, err)
	// }
	// defer f.Close()
	// recipient, err := ReadEntity(f)
	// if err != nil {
	// 	t.Errorf("read entity error: %s\n", err)
	// }
	// pubArmorKey
	recipient, err := ReadEntity(bytes.NewBufferString(pubArmorKey))
	if err != nil {
		t.Errorf("read entity error: %s\n", err)
	}
	f, err := os.Open(fileToEnc)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	dst, err := os.Create(fileToEnc + ".gpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dst.Close()
	err = Encrypt([]*openpgp.Entity{recipient}, nil, f, dst)
	if err != nil {
		t.Errorf("unable to encrypt file:%s err: %s\n", fileToEnc, err)
	}
}

func TestDecrypt(t *testing.T) {

	// pubFilename := "alice.asc"
	fileToDec := "file.txt.gpg"
	file := "file1.txt"
	originalFile := "file.txt"
	passphrase := "2rYiibBZ33f9BMxNE6TEX$H.X4#aPyc2g*(LtRAauiGsn}PJT4{9(J)Xsbe@4Jbr"

	fArmorPriv, err := os.Open(privArmoryFilename)
	if err != nil {
		t.Errorf("unable to open file: %s %s", privFilename, err)
	}
	defer fArmorPriv.Close()
	entityList, err := ReadArmorKeyring(fArmorPriv)
	if err != nil {
		t.Errorf("unable to ReadArmorKeyring %s", privArmoryFilename)
	}

	entity := entityList[0]
	passphraseByte := []byte(passphrase)
	entity.PrivateKey.Decrypt(passphraseByte)
	for _, subkey := range entity.Subkeys {
		err := subkey.PrivateKey.Decrypt(passphraseByte)
		if err != nil {
			t.Errorf("unable to decrypt private key: %s", err)
		}
	}
	f, err := os.Open(fileToDec)
	if err != nil {
		t.Errorf("unable to open file %s, err: %s", fileToDec, err)
	}
	defer f.Close()

	dst, err := os.Create(file)
	if err != nil {
		t.Errorf("unable to open file %s, err: %s", file, err)
	}
	defer dst.Close()

	err = Decrypt(f, dst, entityList)
	if err != nil {
		t.Errorf("unable to decrypt file:%s err: %s\n", fileToDec, err)
	}

	fOrig, err := os.Open(originalFile)
	if err != nil {
		t.Errorf("unable to open file %s, err: %s", originalFile, err)
	}
	_, _ = dst.Seek(0, 0)
	f1, _ := ioutil.ReadAll(fOrig)
	f2, _ := ioutil.ReadAll(dst)
	if !bytes.Equal(f1, f2) {
		t.Errorf("files are not equals %s %s", file, originalFile)
	}

}

// func TestDecryptString(t *testing.T) {

// 	// pubFilename := "alice.asc"
// 	fileToDec := "file.txt.gpg"
// 	file := "file1.txt"
// 	originalFile := "file.txt"
// 	passphrase := "2rYiibBZ33f9BMxNE6TEX$H.X4#aPyc2g*(LtRAauiGsn}PJT4{9(J)Xsbe@4Jbr"

// 	fArmorPriv, err := os.Open(privArmoryFilename)
// 	if err != nil {
// 		t.Errorf("unable to open file: %s %s", privFilename, err)
// 	}
// 	defer fArmorPriv.Close()
// 	entityList, err := ReadArmorKeyring(fArmorPriv)
// 	if err != nil {
// 		t.Errorf("unable to ReadArmorKeyring %s", privArmoryFilename)
// 	}

// 	entity := entityList[0]
// 	passphraseByte := []byte(passphrase)
// 	entity.PrivateKey.Decrypt(passphraseByte)
// 	for _, subkey := range entity.Subkeys {
// 		err := subkey.PrivateKey.Decrypt(passphraseByte)
// 		if err != nil {
// 			t.Errorf("unable to decrypt private key: %s", err)
// 		}
// 	}
// 	f, err := os.Open(fileToDec)
// 	if err != nil {
// 		t.Errorf("unable to open file %s, err: %s", fileToDec, err)
// 	}
// 	defer f.Close()

// 	dst, err := os.Create(file)
// 	if err != nil {
// 		t.Errorf("unable to open file %s, err: %s", file, err)
// 	}
// 	defer dst.Close()

// 	err = Decrypt(f, dst, entityList)
// 	if err != nil {
// 		t.Errorf("unable to decrypt file:%s err: %s\n", fileToDec, err)
// 	}

// 	fOrig, err := os.Open(originalFile)
// 	if err != nil {
// 		t.Errorf("unable to open file %s, err: %s", originalFile, err)
// 	}
// 	_, _ = dst.Seek(0, 0)
// 	f1, _ := ioutil.ReadAll(fOrig)
// 	f2, _ := ioutil.ReadAll(dst)
// 	if !bytes.Equal(f1, f2) {
// 		t.Errorf("files are not equals %s %s", file, originalFile)
// 	}

// }
