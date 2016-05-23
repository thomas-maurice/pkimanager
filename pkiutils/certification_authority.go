package pkiutils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

// Generates a new Certification Authority certificate object
func GenerateNewCertificationAuthority(name pkix.Name, validity int, keyLength int, CRLDistURL string) (*x509.Certificate, *rsa.PrivateKey, error) {
	ca := &x509.Certificate{
		SerialNumber:          big.NewInt(0),
		Subject:               name,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(validity, 0, 0),
		BasicConstraintsValid: true,
		IsCA:        true,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}

	if CRLDistURL != "" {
		ca.CRLDistributionPoints = []string{CRLDistURL}
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return nil, nil, err
	}
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}

	certificate, err := x509.ParseCertificate(caBytes)
	if err != nil {
		return nil, nil, err
	}

	return certificate, privateKey, nil
}
