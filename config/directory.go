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

func getHomeDirectory() (directory string, err error) {
	directory = os.Getenv("HOME")
	if directory != "" {
		return
	}
	if user, err := user.Current(); err == nil {
		directory = user.HomeDir
	}
	return
}

func init() {
	homeDirectory, err := getHomeDirectory()
	if err != nil {
		panic(err)
	}
	ConfigDirectory = filepath.Join(homeDirectory, ".mesg")
	AccountDirectory = filepath.Join(ConfigDirectory, "accounts")
}
