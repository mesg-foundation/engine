package services

import (
	leveldbErrors "github.com/syndtr/goleveldb/leveldb/errors"
)

// NotFound represents a service not found error.
type NotFound struct {
	Hash string
}

func (e NotFound) Error() string {
	return "Database services: Service with hash '" + e.Hash + "' not found"
}

func handleErrorNotFound(err error, hash string) error {
	if err == leveldbErrors.ErrNotFound {
		return NotFound{Hash: hash}
	}
	return err
}
