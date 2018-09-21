package services

import (
	"fmt"

	leveldbErrors "github.com/syndtr/goleveldb/leveldb/errors"
)

// NotFound represents a service not found error.
type NotFound struct {
	Hash string
}

func (e NotFound) Error() string {
	return fmt.Sprintf("Database services: Service %q not found", e.Hash)
}

func handleErrorNotFound(err error, hash string) error {
	if err == leveldbErrors.ErrNotFound {
		return NotFound{Hash: hash}
	}
	return err
}

// DecodeError represents a service impossible to decode
type DecodeError struct {
	Hash string
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("Database services: Could not decode service %q.", e.Hash)
}
