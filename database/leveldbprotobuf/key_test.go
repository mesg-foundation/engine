package leveldbprotobuf

import (
	"testing"

	"github.com/stvp/assert"
)

const testKey = "id_1"

func TestMakeKey(t *testing.T) {
	key := makeKey(testCollection, testKey)
	keyTest := []byte(testCollection + collectionSeparatorKey + testKey)
	assert.Equal(t, keyTest, key)
}

func TestSplitKey(t *testing.T) {
	bytes := makeKey(testCollection, testKey)
	collection, key := splitKey(bytes)
	assert.Equal(t, testCollection, collection)
	assert.Equal(t, testKey, key)
}
