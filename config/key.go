package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// Key is the structure that stores all informations related to the key system
type Key struct {
	Address string
}

// GetKey returns the core's key
func (c *Config) GetKey() (*Key, error) {
	path := filepath.Join(c.Core.Path, "key.yml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := c.createKey(path); err != nil {
			return nil, err
		}
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var key Key
	if err := yaml.UnmarshalStrict(data, &key); err != nil {
		return nil, err
	}
	return &key, nil
}

func (c *Config) createKey(path string) error {
	// TODO: generate ETH address here
	key := &Key{
		Address: "0x2551d2357c8da54b7d330917e0e769d33f1f5b93",
	}
	data, err := yaml.Marshal(key)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, os.ModePerm)
}
