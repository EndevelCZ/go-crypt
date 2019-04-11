package main

/**
  Example hack to encrypt a file using a GPG encryption key. Works with GPG v2.x.
  The encrypted file e.g. /tmp/data.txt.gpg can then be decrypted using the standard command
  gpg /tmp/data.txt.gpg

  Assumes you have **created** an encryption key and exported armored version.
  You have to read the armored key directly as Go cannot read pubring.kbx (yet).

  Export your key using command:
    gpg2 --export --armor [KEY ID] > /tmp/pubKey.asc
*/

import (
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

// change as required
const pubKey = "mykey.asc"
const fileToEnc = "file.txt"

func main() {
	log.Println("Public key:", pubKey)

	// Read in public key
	recipient, err := readEntity(pubKey)
	if err != nil {
		fmt.Println(err)
		return
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
	encrypt([]*openpgp.Entity{recipient}, nil, f, dst)

	os.Open()
}
func encryptSymmetric(w io.Writer, passphrase []byte) error {
	wc, err := openpgp.SymmetricallyEncrypt(w, passphrase, nil, nil)
	if err != nil {
		return err
	}
	defer wc.Close()
	wc.Write([]byte(w))
	// _, err = plaintext.Write(message)
	fmt.Println("encryptSymmetric!")
	// plaintext.Close()
	// w.Close()

	return nil

}
func encrypt(recip []*openpgp.Entity, signer *openpgp.Entity, r io.Reader, w io.Writer) error {
	wc, err := openpgp.Encrypt(w, recip, signer, &openpgp.FileHints{IsBinary: true}, nil)
	fmt.Printf("encrypt %v %T\n", wc, wc)
	if err != nil {
		return err
	}
	if _, err := io.Copy(wc, r); err != nil {
		return err
	}
	return wc.Close()
}

func readEntity(name string) (*openpgp.Entity, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	block, err := armor.Decode(f)
	if err != nil {
		return nil, err
	}
	return openpgp.ReadEntity(packet.NewReader(block.Body))
}
