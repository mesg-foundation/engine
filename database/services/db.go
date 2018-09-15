package services

import (
	"path/filepath"
	"sync"

	"github.com/mesg-foundation/core/config"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	_instance     *leveldb.DB
	instances     = 0
	instanceMutex sync.Mutex
)

func open() (db *leveldb.DB, err error) {
	cfg, err := config.Global()
	if err != nil {
		return nil, err
	}

	instanceMutex.Lock()
	defer instanceMutex.Unlock()
	if _instance == nil {
		storagePath := filepath.Join(cfg.Database.Path, "services")
		_instance, err = leveldb.OpenFile(storagePath, nil)
		if err != nil {
			return nil, err
		}
	}
	instances++
	return _instance, nil
}

func close() error {
	instanceMutex.Lock()
	defer instanceMutex.Unlock()
	if _instance != nil && instances == 0 {
		instances--
		err := _instance.Close()
		_instance = nil
		return err
	}
	return nil
}
