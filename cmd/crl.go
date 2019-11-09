package cmd

import (
	"crypto/rand"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/thomas-maurice/pkimanager/pkiutils"
)

var CRLCmd = &cobra.Command{
	Use:   "crl",
	Short: "Manages CRL",
	Long:  ``,
}

var RegenCRLCmd = &cobra.Command{
	Use:   "regen",
	Short: "Regenerates the CRL",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(path.Join(viper.GetString("ca_root"), "revoked")); err != nil {
			logrus.Info("Creating revoked certificate directory")
			if err := os.MkdirAll(path.Join(viper.GetString("ca_root"), "revoked"), 0750); err != nil {
				logrus.Fatalf("Could not create %s: %s", path.Join(viper.GetString("ca_root"), "revoked"), err)
			}
		}
		CRLFileName := path.Join(viper.GetString("ca_root"), "ca.crl")

		ca, err := pkiutils.LoadCertificate(path.Join(viper.GetString("ca_root"), "ca.crt"))
		if err != nil {
			logrus.Fatalf("Could not load CA certificate: %s", err)
		}
		caKey, err := pkiutils.LoadRSAKey(path.Join(viper.GetString("ca_root"), "ca.key"))
		if err != nil {
			logrus.Fatalf("Could not load CA key: %s", err)
		}

		var revokedCerts []pkix.RevokedCertificate
		revoked, err := ioutil.ReadDir(path.Join(viper.GetString("ca_root"), "revoked"))
		if err != nil {
			logrus.Fatalf("Could not read revoked certificates directory: %s", err)
		}

		for _, certificate := range revoked {
			cert, err := pkiutils.LoadCertificate(path.Join(viper.GetString("ca_root"), "revoked", certificate.Name()))
			if err != nil {
				logrus.Errorf("Could not load CA certificate %s: %s", certificate.Name(), err)
			}
			logrus.Infof("Adding certificate %s, serial %s to revokation list", cert.Subject.CommonName, cert.SerialNumber)
			revokedCerts = append(revokedCerts, pkix.RevokedCertificate{SerialNumber: cert.SerialNumber, RevocationTime: time.Now()})
		}

		crl, err := ca.CreateCRL(rand.Reader, caKey, revokedCerts, time.Now(), time.Now().AddDate(1, 0, 0))
		if err != nil {
			logrus.Fatalf("Could not generate CRL: %s", err)
		}
		pemCrl := pem.EncodeToMemory(&pem.Block{
			Type:  "X509 CRL",
			Bytes: crl,
		})

		if err := ioutil.WriteFile(CRLFileName, pemCrl, 0640); err != nil {
			logrus.Fatalf("Could not write CRL file %s: %s", CRLFileName, err)
		}

		logrus.Infof("Regeneration complete, CRL written to %s", CRLFileName)
	},
}

func InitCRLCmd() {
	CRLCmd.AddCommand(RegenCRLCmd)
}
