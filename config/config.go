package config

import (
	"path"

	"github.com/spf13/viper"
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
