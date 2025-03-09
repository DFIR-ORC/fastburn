package fastfind

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/DFIR-ORC/pkcs7-go"

	log "github.com/sirupsen/logrus"
)

// readRSAPKCS8Key reads a RSA private key in PKCS8 format
func readRSAPKCS8Key(filename string) (*rsa.PrivateKey, error) {
	// Read private key
	privateKeyPEM, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("failed to read private key: %v", err)
		return nil, err
	}

	// Decode PEM to DER
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		log.Tracef("failed to decode private key PEM block")
		return nil, fmt.Errorf("failed to decode private key PEM block")
	}

	// Interpret ASN1 structure
	pkey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Errorf("failed to parse private key: %v", err)
		return nil, err

	}

	// checking if the key is an RSA key
	rsaPrivateKey, ok := pkey.(*rsa.PrivateKey)
	if !ok {
		log.Fatalf("pkey is not an RSA key")
		return nil, fmt.Errorf("pkey is not an RSA key")
	}
	return rsaPrivateKey, nil
}

// DecryptCMSData decrypts a CMS container using a private key
func DecryptCMSData(keyPath string, pkcs7Path string) ([]byte, error) {
	log.Tracef("Decrypting '%s' with key:'%s'", pkcs7Path, keyPath)

	// Read PKCS7 data
	pkcs7Data, err := os.ReadFile(pkcs7Path)
	if err != nil {
		log.Errorf("failed to read pkcs7 data: %v", err)
		return nil, err
	}

	// Parse PKCS7 data
	p7, err := pkcs7.Parse(pkcs7Data)
	if err != nil {
		log.Debugf("Failed to read PKCS7 data from '%s' : %v", pkcs7Path, err)
		return nil, err
	}

	privateKey, err := readRSAPKCS8Key(keyPath)
	if err != nil {
		log.Debugf("Failed to read key from '%s' : %v", keyPath, err)
		return nil, err
	}

	decryptedData, err := p7.DecryptUsingPrivateKey(privateKey)
	if err != nil {
		log.Debugf("Failed to decrypt PKCS7 data from '%s' using '%s' as key: %v", pkcs7Path, keyPath, err)
		return nil, err
	}

	return decryptedData, nil
}

// DecryptCMSContainer decrypts a CMS container using a private key and writes the decrypted data to a file
func DecryptCMSContainer(keyPath string, pkcs7Path string, outPath string) error {

	decryptedData, err := DecryptCMSData(keyPath, pkcs7Path)
	if err != nil {
		return err
	}

	err = os.WriteFile(outPath, decryptedData, 0644)
	if err != nil {
		log.Debugf("Failed to write decrypted data to '%s': %v", outPath, err)
		return err
	}

	return err
}

//eof
