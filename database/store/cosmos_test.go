package store

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store/transient"
	"github.com/stretchr/testify/require"
)

func newCosmosStore(t *testing.T) (*CosmosStore, func()) {
	store := NewCosmosStore(transient.NewStore())
	return store, func() {
		require.NoError(t, store.Close())
	}
}

func TestCosmosStorePut(t *testing.T) {
	store, closer := newCosmosStore(t)
	defer closer()
	require.NoError(t, store.Put([]byte("hello"), []byte("world")))
}

func TestCosmosStoreDelete(t *testing.T) {
	store, closer := newCosmosStore(t)
	defer closer()
	store.Put([]byte("hello"), []byte("world"))
	require.NoError(t, store.Delete([]byte("hello")))
}

func TestCosmosStoreGet(t *testing.T) {
	store, closer := newCosmosStore(t)
	defer closer()
	store.Put([]byte("hello"), []byte("world"))
	value, err := store.Get([]byte("hello"))
	require.NoError(t, err)
	require.Equal(t, []byte("world"), value)
}

func TestCosmosStoreHas(t *testing.T) {
	store, closer := newCosmosStore(t)
	defer closer()
	store.Put([]byte("hello"), []byte("world"))
	has, err := store.Has([]byte("hello"))
	require.NoError(t, err)
	require.True(t, has)
}

func TestCosmosStoreIterate(t *testing.T) {
	store, closer := newCosmosStore(t)
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
	require.Equal(t, -1, i)
	iter.Release()
	require.NoError(t, iter.Error())
}
