package services

import (
	"path/filepath"
	"sync"

	"github.com/mesg-foundation/core/config"
	"github.com/syndtr/goleveldb/leveldb"
)

var _instance *leveldb.DB
var instances = 0
var instanceMutex sync.Mutex

func open() (db *leveldb.DB, err error) {
	instanceMutex.Lock()
	defer instanceMutex.Unlock()
	if _instance == nil {
		storagePath := filepath.Join(config.Path, "database", "services")
		_instance, err = leveldb.OpenFile(storagePath, nil)
		if err != nil {
			return nil, err
		}
	}
	instances++
	db = _instance
	return
}

func close() (err error) {
	instanceMutex.Lock()
	defer instanceMutex.Unlock()
	instances--
	if _instance != nil && instances == 0 {
		err = _instance.Close() // TODO: this should return the error before set _instance to nil
		_instance = nil
	}
	return
}
