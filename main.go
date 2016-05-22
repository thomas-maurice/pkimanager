package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/thomas-maurice/pkimanager/cmd"
	"github.com/thomas-maurice/pkimanager/config"
)

func main() {
	cmd.InitRootCmd()
	if err := config.InitConfig(); err != nil {
		logrus.Errorf("Could not load configuration file: %s", err)
	}
	if err := cmd.RootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
