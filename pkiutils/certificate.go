package pkiutils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/satori/go.uuid"
	"math/big"
	"strings"
	"time"
)

// Generates a new Certificate object
func GenerateNewCertificate(name pkix.Name, validity int, keyLength int, ca *x509.Certificate, caKey *rsa.PrivateKey, clientCert bool, altnames []string) (*x509.Certificate, *rsa.PrivateKey, error) {
	number := uuid.NewV4().Bytes()
	serial := big.NewInt(0)
	serial.SetBytes(number)
	certificate := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               name,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(validity, 0, 0),
		BasicConstraintsValid: true,
		IsCA:     false,
		KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageKeyAgreement,
	}

	if len(altnames) != 0 {
		for _, altname := range altnames {
			if strings.TrimSpace(altname) != "" {
				certificate.DNSNames = append(certificate.DNSNames, strings.TrimSpace(altname))
			}
		}
	}

	if clientCert == true {
		certificate.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	} else {
		certificate.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return nil, nil, err
	}
	certificateBytes, err := x509.CreateCertificate(rand.Reader, certificate, ca, &privateKey.PublicKey, caKey)
	if err != nil {
		return nil, nil, err
	}

	certif, err := x509.ParseCertificate(certificateBytes)
	if err != nil {
		return nil, nil, err
	}

	return certif, privateKey, nil
}
