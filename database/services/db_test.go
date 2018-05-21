package services

import (
	"testing"

	"github.com/stvp/assert"
)

func TestDb(t *testing.T) {
	db := open()
	assert.NotNil(t, db)
	close()
}
