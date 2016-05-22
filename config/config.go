package config

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path"
)

func InitConfig() error {
	viper.SetConfigName("pkimanager")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/")

	viper.SetDefault("ca_root", "./ca_root")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	if _, err := os.Stat(viper.GetString("ca_root")); err != nil {
		logrus.Infof("Creating directory %s", viper.Get("ca_root"))
		if err := os.MkdirAll(viper.GetString("ca_root"), 0750); err != nil {
			logrus.Fatalf("Could not create %s: %s", viper.Get("ca_root"), err)
		}
	}

	return nil
}

// Returns the path to the revoked certificates directory
func GetRevokedPath() string {
	return path.Join(viper.GetString("ca_root"), "revoked")
}

// Returns the path to the certificates directory
func GetCertificatePath() string {
	return path.Join(viper.GetString("ca_root"), "certificates")
}

// Returns the path to the keys directory
func GetKeysPath() string {
	return path.Join(viper.GetString("ca_root"), "keys")
}
