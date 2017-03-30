package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thomas-maurice/pkimanager/cmd"
	"github.com/thomas-maurice/pkimanager/config"
)

func main() {
	cmd.InitRootCmd()
	if err := config.InitConfig(); err != nil {
		logrus.WithError(err).Warning("Could not load configuration file")
	}

	if _, err := os.Stat(viper.GetString("ca_root")); err != nil {
		logrus.Infof("Creating directory %s", viper.Get("ca_root"))
		if err := os.MkdirAll(viper.GetString("ca_root"), 0750); err != nil {
			logrus.Fatalf("Could not create %s: %s", viper.Get("ca_root"), err)
		}
	}

	if err := cmd.RootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
