package main

import (
	"flag"
	"fmt"

	"github.com/EndevelCZ/go-crypt/crypto"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	plaintextFile := flag.String("plaintextFile", "file.txt", "text file do encrypt (path to file)")
	ciphertextFile := flag.String("ciphertextFile", "file.txt.gpg", "encrypted output file (path to file)")
	passphrase := flag.String("passphrase", "password", "encryption passphrase")
	err := crypto.GpgEncryptFileSymmetric(plaintextFile, ciphertextFile, passphrase)
	if err != nil {
		log.fatal(err)
	} else {
		fmt.Printf("%s has been saved and encrypted with passphrase: %s", ciphertextFile, passphrase)
	}
	// f, err := os.Open("file.txt")
	// check(err)
	// encryptionPassphrase := []byte("heslo")
	// plaintext, err := openpgp.SymmetricallyEncrypt(f, encryptionPassphrase, nil, nil)

	// defer plaintext.Close()

}

// encryptionPassphrase := []byte("golang")
// encryptionText := "Hello world. Encryption and Decryption testing.\n"
// encryptionType := "PGP SIGNATURE"

// encbuf := bytes.NewBuffer(nil)
// w, err := armor.Encode(encbuf, encryptionType, nil)
// if err != nil {
// 	log.Fatal(err)
// }

// plaintext, err := openpgp.SymmetricallyEncrypt(w, encryptionPassphrase, nil, nil)
// if err != nil {
// 	log.Fatal(err)
// }
// message := []byte(encryptionText)
// _, err = plaintext.Write(message)

// plaintext.Close()
// w.Close()

// // permissions := 0644 // or whatever you need
// // permissions := int(0777)

// err = ioutil.WriteFile("file.txt.gpg", []byte(encbuf.String()), 0777)
// if err != nil {
// 	// handle error
// 	// check(err)
// }

// fmt.Printf("Encrypted:\n%s\n", encbuf)

// decbuf := bytes.NewBuffer([]byte(encbuf.String()))
// result, err := armor.Decode(decbuf)
// if err != nil {
// 	log.Fatal(err)
// }

// md, err := openpgp.ReadMessage(result.Body, nil, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
// 	return encryptionPassphrase, nil
// }, nil)
// if err != nil {
// 	log.Fatal(err)
// }

// bytes, err := ioutil.ReadAll(md.UnverifiedBody)
// fmt.Printf("Decrypted:\n%s\n", string(bytes))

// package main

// import (
// 	"fmt"
// 	"io/ioutil"

// 	"github.com/EndevelCZ/go-crypt/crypto"
// )

// func check(e error) {
// 	if e != nil {
// 		panic(e)
// 	}
// }
// func main() {
// 	// kmsKey := flag.String("kms-key-path", "nil", "path for kms key")

// 	// flag.Parse()
// 	// fmt.Println("kms-key-path:", *kmsKey)
// 	// os.Exit(3)
// 	// kmsKey := "projects/cterminal-kms/locations/europe-west3/keyRings/cterminal-store/cryptoKeys/store-test-dev-symetric"
// 	// plaintext, err := ioutil.ReadFile("README.md")
// 	// check(err)
// 	// ciphertext, err := crypto.KmsEncryptSymmetric(kmsKey, plaintext)
// 	// check(err)
// 	// fmt.Println(string(ciphertext))
// 	// fmt.Println("-----------------")
// 	// // plaintext, err = crypto.Decrypt #.DecryptSymmetric(key, ciphertext)
// 	// plaintext, err = crypto.KmsDecryptSymmetric(kmsKey, ciphertext)
// 	// check(err)
// 	// fmt.Print(string(plaintext))

// 	plaintext, err := ioutil.ReadFile("README.md")
// 	w := bufio.NewWriter(f)
// 	check(err)
// 	passphrase := []byte("password")
// 	ciphertext, err := crypto.GpgEncryptSymmetric(plaintext, passphrase)
// 	check(err)
// 	fmt.Println(string(ciphertext))

// }
