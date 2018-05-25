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
var instanceMutex sync.Mutex
var databaseLock sync.WaitGroup

func open() (db *leveldb.DB, err error) {
	instanceMutex.Lock()
	defer instanceMutex.Unlock()
	databaseLock.Wait()
	databaseLock.Add(1)
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
	instanceMutex.Lock()
	defer instanceMutex.Unlock()
	defer databaseLock.Done()
	if _instance != nil {
		err = _instance.Close()
		_instance = nil
	}
	return
}
