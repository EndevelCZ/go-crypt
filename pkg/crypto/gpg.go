package crypto

import (
	"io"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

// ReadEntity reads pub armored key
func ReadEntity(r io.Reader) (*openpgp.Entity, error) {
	block, err := armor.Decode(r)
	if err != nil {
		return nil, err
	}
	return openpgp.ReadEntity(packet.NewReader(block.Body))
}

func ReadEntities(r ...io.Reader) ([]*openpgp.Entity, error) {
	entities := []*openpgp.Entity{}
	for _, rr := range r {
		block, err := armor.Decode(rr)
		if err != nil {
			return nil, err
		}
		entity, err := openpgp.ReadEntity(packet.NewReader(block.Body))
		if err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}
	return entities, nil
}

func ReadPrivateKey(r io.Reader) (*openpgp.Entity, error) {
	return openpgp.ReadEntity(packet.NewReader(r))
}
func ReadKeyring(r io.Reader) (openpgp.EntityList, error) {
	return openpgp.ReadKeyRing(r)
}
func ReadArmorKeyring(r io.Reader) (openpgp.EntityList, error) {
	return openpgp.ReadArmoredKeyRing(r)
}

func Encrypt(recip []*openpgp.Entity, signer *openpgp.Entity, r io.Reader, w io.Writer) error {
	wc, err := openpgp.Encrypt(w, recip, signer, &openpgp.FileHints{IsBinary: true}, nil)
	if err != nil {
		return err
	}
	if _, err := io.Copy(wc, r); err != nil {
		return err
	}
	return wc.Close()
}
func Decrypt(r io.Reader, w io.Writer, ent openpgp.EntityList) error {
	md, err := openpgp.ReadMessage(r, ent, nil, nil)
	if err != nil {
		return err
	}
	if _, err := io.Copy(w, md.UnverifiedBody); err != nil {
		return err
	}
	return nil
}

// GpgEncryptFileSymmetric open plaintextFile and encrypt it with symmetric
// gpg cypher. The result is stored in ciphertextFile
func GpgEncryptSymmetric(fp io.Reader, fc io.Writer, passphrase []byte) error {

	plaintext, err := openpgp.SymmetricallyEncrypt(fc, passphrase, nil, nil)
	if err != nil {
		return err
	}
	defer plaintext.Close()
	if _, err := io.Copy(plaintext, fp); err != nil {
		return err
	}
	return nil
}
func GpgDecryptSymmetric(fc io.Reader, fp io.Writer, passphrase []byte) error {
	md, err := openpgp.ReadMessage(fc, nil, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		return passphrase, nil
	}, nil)
	if err != nil {
		return err
	}
	if _, err := io.Copy(fp, md.UnverifiedBody); err != nil {
		return err
	}
	return nil
}
