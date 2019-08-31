package store

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const testLevelDBPath = "test"

func newLevelDBStore(t *testing.T) (*LevelDBStore, func()) {
	store, err := NewLevelDBStore(testLevelDBPath)
	require.NoError(t, err)
	return store, func() {
		require.NoError(t, store.Close())
		os.RemoveAll(testLevelDBPath)
	}
}

func TestLevelDBStorePut(t *testing.T) {
	store, closer := newLevelDBStore(t)
	defer closer()
	require.NoError(t, store.Put([]byte("hello"), []byte("world")))
}

func TestLevelDBStoreDelete(t *testing.T) {
	store, closer := newLevelDBStore(t)
	defer closer()
	store.Put([]byte("hello"), []byte("world"))
	require.NoError(t, store.Delete([]byte("hello")))
}

func TestLevelDBStoreGet(t *testing.T) {
	store, closer := newLevelDBStore(t)
	defer closer()
	store.Put([]byte("hello"), []byte("world"))
	value, err := store.Get([]byte("hello"))
	require.NoError(t, err)
	require.Equal(t, []byte("world"), value)
}

func TestLevelDBStoreHas(t *testing.T) {
	store, closer := newLevelDBStore(t)
	defer closer()
	store.Put([]byte("hello"), []byte("world"))
	has, err := store.Has([]byte("hello"))
	require.NoError(t, err)
	require.True(t, has)
}

func TestLevelDBStoreIterate(t *testing.T) {
	store, closer := newLevelDBStore(t)
	defer closer()

	data := []struct {
		key   []byte
		value []byte
	}{
		{key: []byte("hello"), value: []byte("world")},
		{key: []byte("foo"), value: []byte("bar")},
	}
	for _, d := range data {
		store.Put(d.key, d.value)
	}
	iter := store.NewIterator()
	i := len(data) - 1
	for iter.Next() {
		require.Equal(t, data[i].key, iter.Key())
		require.Equal(t, data[i].value, iter.Value())
		i--
	}
	iter.Release()
	require.NoError(t, iter.Error())
}
