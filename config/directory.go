package config

import (
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

const configDirectory = ".mesg"

func getConfigPath() (string, error) {
	homePath, err := homedir.Dir()
	if err != nil {
		return "", nil
	}
	return filepath.Join(homePath, configDirectory), nil
}

func createConfigPath() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}
	err = os.Mkdir(configPath, os.ModePerm)
	if os.IsExist(err) == false {
		return err
	}
	return nil
}
