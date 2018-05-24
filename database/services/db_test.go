package services

import (
	"testing"

	"github.com/stvp/assert"
)

func TestDb(t *testing.T) {
	db, err := open()
	assert.Nil(t, err)
	assert.NotNil(t, db)
	close()
}
