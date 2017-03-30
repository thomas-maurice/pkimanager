package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configfile string

var RootCmd = &cobra.Command{
	Use:   "camanager",
	Short: "Certification Authority Manager",
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Certification Authority Manager v0.3")
	},
}

func InitRootCmd() {
	InitCACmd()
	InitCRLCmd()
	InitCertificateCmd()
	RootCmd.AddCommand(VersionCmd)
	RootCmd.AddCommand(CACmd)
	RootCmd.AddCommand(CRLCmd)
	RootCmd.AddCommand(CertificateCmd)
}
