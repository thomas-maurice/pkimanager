package cmd

import (
	"crypto/x509/pkix"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thomas-maurice/pkimanager/pkiutils"
	"os"
	"path"
)

var forceReplaceCA bool
var CAValidity int
var CAKeySize int
var CACountry string
var CALocality string
var CAProvince string
var CAOrganization string
var CAOrganizationalUnit string
var CACRLDistURL string

var CACmd = &cobra.Command{
	Use:   "ca",
	Short: "Manages a certification authority",
	Long:  ``,
}

var CreateCACmd = &cobra.Command{
	Use:   "create [CommonName]",
	Short: "Creates a new certification authority",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		certFileName := path.Join(viper.GetString("ca_root"), "ca.crt")
		keyFileName := path.Join(viper.GetString("ca_root"), "ca.key")

		if _, err := os.Stat(certFileName); err == nil {
			if forceReplaceCA {
				logrus.Warningf("The file %s already exists and will be overwritten", certFileName)
			} else {
				logrus.Fatalf("The file %s already exists, use --force to overwrite it", certFileName)
			}
		}

		if len(args) != 1 {
			logrus.Fatal("The CommonName is mandatory")
		}
		name := pkix.Name{
			CommonName: args[0],
		}

		if CACountry != "" {
			name.Country = []string{CACountry}
		}

		if CALocality != "" {
			name.Locality = []string{CALocality}
		}

		if CAProvince != "" {
			name.Province = []string{CAProvince}
		}

		if CAOrganization != "" {
			name.Organization = []string{CAOrganization}
		}

		if CAOrganizationalUnit != "" {
			name.OrganizationalUnit = []string{CAOrganizationalUnit}
		}

		logrus.Infof("Generating the CA %s with a %d bits key and a %d years validity", args[0], CAKeySize, CAValidity)
		ca, key, err := pkiutils.GenerateNewCertificationAuthority(name, CAValidity, CAKeySize, CACRLDistURL)
		if err != nil {
			logrus.Fatal(err)
		}
		if err := pkiutils.WriteCertificate(certFileName, ca); err != nil {
			logrus.Fatal(err)
		}
		if err := pkiutils.WriteRSAKey(keyFileName, key); err != nil {
			logrus.Fatal(err)
		}

		logrus.Infof("Generation complete, certificate written to %s, and key to %s", certFileName, keyFileName)
	},
}

var ShowCACmd = &cobra.Command{
	Use:   "show",
	Short: "Shows the current CA",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		certFileName := path.Join(viper.GetString("ca_root"), "ca.crt")
		//keyFileName := path.Join(viper.GetString("ca_root"), "ca.key")

		if _, err := os.Stat(certFileName); err != nil {
			logrus.Fatalf("Could not state file %s: %s", certFileName, err)
		}

		ca, err := pkiutils.LoadCertificate(path.Join(viper.GetString("ca_root"), "ca.crt"))
		if err != nil {
			logrus.Fatalf("Could not load CA certificate: %s", err)
		}

		logrus.Infof("Common Name: %s, Serial: %s", ca.Subject.CommonName, ca.SerialNumber)
		logrus.Infof("Country: %s, Locality: %s", ca.Subject.Country, ca.Subject.Locality)
		logrus.Infof("Expiry: %s", ca.NotAfter)
	},
}

func InitCACmd() {
	CreateCACmd.Flags().BoolVarP(&forceReplaceCA, "force", "f", false, "Force the replacement of existing files")
	CreateCACmd.Flags().IntVarP(&CAValidity, "validity", "v", 1, "For how many years will this CA be valid")
	CreateCACmd.Flags().IntVarP(&CAKeySize, "key-size", "k", 4096, "CA private key size")
	CreateCACmd.Flags().StringVarP(&CACountry, "country", "c", "", "CA country code")
	CreateCACmd.Flags().StringVarP(&CAOrganization, "org", "o", "", "CA organization")
	CreateCACmd.Flags().StringVarP(&CAOrganizationalUnit, "ou", "", "", "CA organizational unit")
	CreateCACmd.Flags().StringVarP(&CAProvince, "province", "p", "", "CA province")
	CreateCACmd.Flags().StringVarP(&CALocality, "locality", "l", "", "CA locality")
	CreateCACmd.Flags().StringVarP(&CACRLDistURL, "crl-dist-url", "", "", "CA Certificate Revocation List distribution URL")
	CACmd.AddCommand(CreateCACmd)
	CACmd.AddCommand(ShowCACmd)
}
