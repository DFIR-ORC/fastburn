package fastfind

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"go.mozilla.org/pkcs7"
)

func decryptPKCS7Data(pkcs7Data, privateKeyData, certificateData []byte) ([]byte, error) {
	log.Trace("Parsing PKCS7 data")
	p7, err := pkcs7.Parse(pkcs7Data)
	if err != nil {
		return nil, err
	}

	log.Trace("Decoding private key")
	blockKey, _ := pem.Decode(privateKeyData)
	if blockKey == nil || blockKey.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM private key")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(blockKey.Bytes)
	if err != nil {
		return nil, err
	}

	log.Trace("Decoding certificate")
	blockCert, _ := pem.Decode(certificateData)
	if blockCert == nil || blockCert.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("failed to decode PEM certificate")
	}
	certificate, err := x509.ParseCertificate(blockCert.Bytes)
	if err != nil {
		return nil, err
	}

	log.Trace("Decrypting")
	decryptedData, err := p7.Decrypt(certificate, privateKey)
	if err != nil {
		return nil, err
	}

	log.Trace("PKCS7 decrypted")
	return decryptedData, nil
}

func DecryptPKCSData(certPath string, keyPath string, pkcs7Path string) ([]byte, error) {
	pkcs7Data, err := os.ReadFile(pkcs7Path)
	if err != nil {
		log.Warnf("Failed to read PKCS7 data from '%s' : %v", pkcs7Path, err)
		return nil, err
	}

	certificateData, err := os.ReadFile(certPath)
	if err != nil {
		log.Warnf("Failed to read certificate from '%s' : %v", certPath, err)
		return nil, err
	}

	privateKeyData, err := os.ReadFile(keyPath)
	if err != nil {
		log.Warnf("Failed to read key from '%s' : %v", keyPath, err)
		return nil, err
	}

	decryptedData, err := decryptPKCS7Data(pkcs7Data, privateKeyData, certificateData)
	if err != nil {
		log.Warnf("Failed to decrypt PKCS7 data from '%s' using '%s' as certificate and '%s' as key: %v", pkcs7Path, certPath, keyPath, err)
		return nil, err
	}

	return decryptedData, nil
}

func DecryptPKCS7Container(certPath string, keyPath string, pkcs7Path string, outPath string) error {

	decryptedData, err := DecryptPKCSData(certPath, keyPath, pkcs7Path)
	if err != nil {
		return err
	}

	err = os.WriteFile(outPath, decryptedData, 0644)
	if err != nil {
		log.Warnf("Failed to write decrypted data to '%s': %v", outPath, err)
		return err
	}

	return err
}
