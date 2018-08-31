package config

import (
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	defaultDirectory        = ".mesg"
	defaultServiceDirectory = "services"
	MESGPath                = "MESG.Path"
)

func setDirectoryDefault() {
	viper.SetDefault(MESGPath, getDefaultPath())
}

func getDefaultPath() string {
	path, err := homedir.Dir()
	if err != nil {
		return ""
	}
	return filepath.Join(path, defaultDirectory)
}

func createPath() error {
	path := viper.GetString(MESGPath)
	err := os.Mkdir(path, os.ModePerm)
	if os.IsExist(err) == false {
		return err
	}
	return nil
}

func createServicesPath() error {
	path := filepath.Join(viper.GetString(MESGPath), defaultServiceDirectory)
	err := os.Mkdir(path, os.ModePerm)
	if os.IsExist(err) == false {
		return err
	}
	return nil
}
