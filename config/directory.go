package config

import (
	"os"
	"os/user"
	"path/filepath"
)

// ConfigDirectory is the root directory where everything related to MESG is stored
var ConfigDirectory string

// AccountDirectory is the directory where all accounts are stored
var AccountDirectory string

func detectHomePath() (path string, err error) {
	user, err := user.Current()
	if err != nil {
		return
	}
	path = user.HomeDir
	return
}

func getHomeDirectory() (directory string, err error) {
	directory = os.Getenv("HOME")
	if directory != "" {
		return
	}
	directory, err = detectHomePath()
	return
}

func getConfigDirectory() (directory string, err error) {
	homeDirectory, err := getHomeDirectory()
	if err != nil {
		panic(err)
	}
	directory = filepath.Join(homeDirectory, ".mesg")
	return
}

func getAccountDirectory() (directory string, err error) {
	configDirectory, err := getConfigDirectory()
	directory = filepath.Join(configDirectory, "accounts")
	return
}

func init() {
	var err error
	ConfigDirectory, err = getConfigDirectory()
	AccountDirectory, err = getAccountDirectory()
	os.Mkdir(ConfigDirectory, os.ModePerm)
	os.Mkdir(AccountDirectory, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
