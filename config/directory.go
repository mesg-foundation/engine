package config

import (
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

const configDirectory = ".mesg"

func getConfigDirectory() (string, error) {
	homeDirectory, err := homedir.Dir()
	if err != nil {
		return "", nil
	}
	return filepath.Join(homeDirectory, configDirectory), nil
}

func createConfigDirectory() error {
	configDirectory, err := getConfigDirectory()
	if err != nil {
		return err
	}
	err = os.Mkdir(configDirectory, os.ModePerm)
	if os.IsExist(err) == false {
		return err
	}
	return nil
}

func init() {
	err := createConfigDirectory()
	if err != nil {
		panic(err)
	}
}
