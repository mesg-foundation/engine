package services

import (
	"path/filepath"
	"sync"

	"github.com/spf13/viper"

	"github.com/mesg-foundation/core/config"
	"github.com/syndtr/goleveldb/leveldb"
)

var storagePath = filepath.Join(viper.GetString(config.MESGPath), "database", "services")
var _instance *leveldb.DB
var mu sync.Mutex

func open() (db *leveldb.DB, err error) {
	mu.Lock()
	defer mu.Unlock()
	if _instance == nil {
		_instance, err = leveldb.OpenFile(storagePath, nil)
		if err != nil {
			panic(err)
		}
	}
	db = _instance
	return
}

func close() (err error) {
	mu.Lock()
	defer mu.Unlock()
	if _instance != nil {
		err = _instance.Close()
		_instance = nil
	}
	return
}
