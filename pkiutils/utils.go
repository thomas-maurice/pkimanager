package pkiutils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

// Writes a certificate to a file
func WriteCertificate(filename string, certificate *x509.Certificate) error {
	block := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certificate.Raw,
	})
	return ioutil.WriteFile(filename, block, 0640)
}

// Loads a certificate from a file
func LoadCertificate(filename string) (*x509.Certificate, error) {
	encodedBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	certificateBytes, _ := pem.Decode(encodedBytes)
	certificate, err := x509.ParseCertificate(certificateBytes.Bytes)
	if err != nil {
		return nil, err
	}
	return certificate, nil
}

// Writes an RSA private key to a file
func WriteRSAKey(filename string, key *rsa.PrivateKey) error {
	privateBytes := x509.MarshalPKCS1PrivateKey(key)
	encodedBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateBytes},
	)
	return ioutil.WriteFile(filename, encodedBytes, 0640)
}

// Loads an RSA key from a file
func LoadRSAKey(filename string) (*rsa.PrivateKey, error) {
	encodedBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(encodedBytes)
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
