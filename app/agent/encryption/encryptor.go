package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type Encryptor struct {
	rsaPub *rsa.PublicKey
	path   string
}

func NewEncryptor(path string) (*Encryptor, error) {
	encryptor := &Encryptor{}
	err := encryptor.loadPublicKey(path)
	if err != nil {
		return nil, err
	}
	return encryptor, nil
}

// loadPublicKey create a function to load the public key from a file:
func (e *Encryptor) loadPublicKey(path string) error {
	e.path = path
	// if empty path, then no encryption
	if path == "" {
		return nil
	}

	keyData, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return fmt.Errorf("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("not RSA public key")
	}
	e.rsaPub = rsaPub
	return nil
}

// EncryptMessage encrypts the message using the public key:
func (e *Encryptor) EncryptMessage(message string) (string, error) {
	// if empty path, then no encryption
	if e.path == "" {
		return message, nil
	}
	encryptedMessage, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		e.rsaPub,
		[]byte(message),
		nil,
	)
	return string(encryptedMessage), err
}
