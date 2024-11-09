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

type Decryptor struct {
	rsaPriv *rsa.PrivateKey
}

func NewDecryptor(path string) (*Decryptor, error) {
	decryptor := &Decryptor{}
	err := decryptor.loadPrivateKey(path)
	if err != nil {
		return nil, err
	}
	return decryptor, nil
}

// loadPrivateKey loads the private key from a file
func (d *Decryptor) loadPrivateKey(path string) error {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "PRIVATE KEY" {
		return fmt.Errorf("failed to decode PEM block containing private key")
	}

	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	rsaPriv, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return fmt.Errorf("not RSA private key")
	}
	d.rsaPriv = rsaPriv
	return nil
}

// DecryptMessage decrypts the message using the private key
func (d *Decryptor) DecryptMessage(encryptedMessage string) (string, error) {
	decodedMessage := []byte(encryptedMessage)
	decryptedMessage, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		d.rsaPriv,
		decodedMessage,
		nil,
	)
	if err != nil {
		return "", err
	}
	return string(decryptedMessage), nil
}
