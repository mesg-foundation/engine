package services

import (
	"testing"

	"github.com/stvp/assert"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestDb(t *testing.T) {
	db, err := open()
	assert.Nil(t, err)
	assert.NotNil(t, db)
	close()
}

func TestOpenError(t *testing.T) {
	dbx, _ := leveldb.OpenFile(storagePath, nil)
	db, err := open()
	assert.Equal(t, err.Error(), "resource temporarily unavailable")
	assert.Nil(t, db)
	dbx.Close()
	close()
}
