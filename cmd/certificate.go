package cmd

import (
	"crypto/x509/pkix"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thomas-maurice/pkimanager/config"
	"github.com/thomas-maurice/pkimanager/pkiutils"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var forceReplaceCertificate bool
var copyCASubject bool
var ClientCertificate bool
var CertificateValidity int
var CertificateKeySize int
var CertificateCountry string
var CertificateLocality string
var CertificateProvince string
var CertificateOrganization string
var CertificateOrganizationalUnit string
var CertificateAlternateNames string
var CertificateIPs string

var CertificateCmd = &cobra.Command{
	Use:   "certificate",
	Short: "Manages certificates",
	Long:  ``,
}

var CreateCertificateCmd = &cobra.Command{
	Use:   "create [CommonName]",
	Short: "Creates a new certificate",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			logrus.Fatal("The CommonName is mandatory")
		}

		if _, err := os.Stat(config.GetCertificatePath()); err != nil {
			logrus.Info("Creating certificates directory")
			if err := os.MkdirAll(config.GetCertificatePath(), 0750); err != nil {
				logrus.Fatalf("Could not create %s: %s", config.GetCertificatePath(), err)
			}
		}

		if _, err := os.Stat(config.GetKeysPath()); err != nil {
			logrus.Info("Creating keys directory")
			if err := os.MkdirAll(config.GetKeysPath(), 0750); err != nil {
				logrus.Fatalf("Could not create %s: %s", config.GetKeysPath(), err)
			}
		}

		certFileName := path.Join(config.GetCertificatePath(), args[0]+".crt")
		keyFileName := path.Join(config.GetKeysPath(), args[0]+".key")

		if _, err := os.Stat(certFileName); err == nil {
			if forceReplaceCertificate {
				logrus.Warningf("The file %s already exists and will be overwritten", certFileName)
			} else {
				logrus.Fatalf("The file %s already exists, use --force to overwrite it", certFileName)
			}
		}

		ca, err := pkiutils.LoadCertificate(path.Join(viper.GetString("ca_root"), "ca.crt"))
		if err != nil {
			logrus.Fatalf("Could not load CA certificate: %s", err)
		}
		caKey, err := pkiutils.LoadRSAKey(path.Join(viper.GetString("ca_root"), "ca.key"))
		if err != nil {
			logrus.Fatalf("Could not load CA key: %s", err)
		}

		name := pkix.Name{}

		if copyCASubject {
			logrus.Info("Copying Certification Authority's subject")
			name = ca.Subject
		}

		name.CommonName = args[0]

		if CertificateCountry != "" {
			name.Country = []string{CertificateCountry}
		}

		if CertificateLocality != "" {
			name.Locality = []string{CertificateLocality}
		}

		if CertificateProvince != "" {
			name.Province = []string{CertificateProvince}
		}

		if CertificateOrganization != "" {
			name.Organization = []string{CertificateOrganization}
		}

		if CertificateOrganizationalUnit != "" {
			name.OrganizationalUnit = []string{CertificateOrganizationalUnit}
		}

		logrus.Infof("Generating the certificate %s with a %d bits key and a %d years validity", args[0], CertificateKeySize, CertificateValidity)
		certificate, key, err := pkiutils.GenerateNewCertificate(name,
			CertificateValidity,
			CertificateKeySize,
			ca,
			caKey,
			ClientCertificate,
			strings.Split(CertificateAlternateNames, ","),
			strings.Split(CertificateIPs, ","),
		)

		if err != nil {
			logrus.Fatalf("Could not generate certificate: %s", err)
		}
		if err := pkiutils.WriteCertificate(certFileName, certificate); err != nil {
			logrus.Fatalf("Could not write certificate: %s", err)
		}
		if err := pkiutils.WriteRSAKey(keyFileName, key); err != nil {
			logrus.Fatalf("Could not write private key: %s", err)
		}

		logrus.Infof("Generation complete, certificate written to %s, and key to %s", certFileName, keyFileName)
	},
}

var ListCertificateCmd = &cobra.Command{
	Use:   "list",
	Short: "List the certificates",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(config.GetCertificatePath()); err != nil {
			logrus.Fatalf("There is not certificate directory")
		}

		certificates, err := ioutil.ReadDir(config.GetCertificatePath())
		if err != nil {
			logrus.Fatalf("Could not read certificate directory: %s", err)
		}

		for _, certificate := range certificates {
			cert, err := pkiutils.LoadCertificate(path.Join(config.GetCertificatePath(), certificate.Name()))
			if err != nil {
				logrus.Errorf("Could not load CA certificate %s: %s", certificate.Name(), err)
			}
			fmt.Println(cert.Subject.CommonName)
		}
	},
}

var RevokeCertificateCmd = &cobra.Command{
	Use:   "revoke [CommonName]",
	Short: "Revoke a certificates",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			logrus.Fatal("The CommonName is mandatory")
		}

		if _, err := os.Stat(config.GetRevokedPath()); err != nil {
			logrus.Info("Creating revoked certificates directory")
			if err := os.MkdirAll(config.GetRevokedPath(), 0750); err != nil {
				logrus.Fatalf("Could not create %s: %s", config.GetRevokedPath(), err)
			}
		}

		if _, err := os.Stat(path.Join(config.GetCertificatePath(), fmt.Sprintf("%s.crt", args[0]))); err != nil {
			logrus.Fatalf("File %s does not exist", path.Join(config.GetRevokedPath(), fmt.Sprintf("%s.crt", args[0])))
		}
		cert, err := pkiutils.LoadCertificate(path.Join(config.GetCertificatePath(), fmt.Sprintf("%s.crt", args[0])))
		if err != nil {
			logrus.Errorf("Could not load certificate %s: %s", args[0], err)
		}

		err = os.Rename(
			path.Join(config.GetCertificatePath(), fmt.Sprintf("%s.crt", args[0])),
			path.Join(config.GetRevokedPath(), fmt.Sprintf("%s-%s.crt", args[0], cert.SerialNumber)),
		)
		if err != nil {
			logrus.Fatalf("Could not move certificate to revoked directory: %s", err)
		}

		logrus.Info("Certificate revoked, run the crl regen command to make it effective")
	},
}

func InitCertificateCmd() {
	CreateCertificateCmd.Flags().BoolVarP(&forceReplaceCertificate, "force", "f", false, "Force the replacement of existing files")
	CreateCertificateCmd.Flags().StringVarP(&CertificateAlternateNames, "alternate-names", "a", "", "Alternate names for the certificate, coma separated")
	CreateCertificateCmd.Flags().StringVarP(&CertificateIPs, "ips", "i", "", "IP addresses for the certificate, coma separated")
	CreateCertificateCmd.Flags().BoolVarP(&ClientCertificate, "client", "", false, "Is the certificate designed for a client")
	CreateCertificateCmd.Flags().BoolVarP(&copyCASubject, "ca-subject", "", false, "Copy CA Subject, except for CommonName")
	CreateCertificateCmd.Flags().IntVarP(&CertificateValidity, "validity", "v", 1, "For how many years will this certificate be valid")
	CreateCertificateCmd.Flags().IntVarP(&CertificateKeySize, "key-size", "k", 4096, "Certificate key size, a high number is prefered")
	CreateCertificateCmd.Flags().StringVarP(&CertificateCountry, "country", "c", "", "Certificate country code")
	CreateCertificateCmd.Flags().StringVarP(&CertificateOrganization, "org", "o", "", "Certificate organization")
	CreateCertificateCmd.Flags().StringVarP(&CertificateOrganizationalUnit, "ou", "", "", "Certificate organizational unit")
	CreateCertificateCmd.Flags().StringVarP(&CertificateProvince, "province", "p", "", "Certificate province")
	CreateCertificateCmd.Flags().StringVarP(&CertificateLocality, "locality", "l", "", "Certificate locality")
	CertificateCmd.AddCommand(ListCertificateCmd)
	CertificateCmd.AddCommand(CreateCertificateCmd)
	CertificateCmd.AddCommand(RevokeCertificateCmd)
}
