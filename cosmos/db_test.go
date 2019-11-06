package cosmos

import (
	"errors"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/transient"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/hash"
	"github.com/stretchr/testify/require"
)

type TestDBData struct {
	A string
	B string
}

func TestCosmosDbTyped(t *testing.T) {
	cdc := codec.New()
	db := NewDB(transient.NewStore(), cdc)
	cdc.RegisterConcrete(&TestDBData{}, "TestDBData", nil)
	var (
		key  = hash.Int(1)
		data = TestDBData{
			A: "test",
		}
		key2  = hash.Int(2)
		data2 = TestDBData{
			A: "foo",
		}
	)
	t.Run("Save", func(t *testing.T) {
		require.NoError(t, db.Save(key, data))
		require.NoError(t, db.Save(key2, data2))
	})
	t.Run("Get", func(t *testing.T) {
		var d TestDBData
		require.NoError(t, db.Get(key, &d))
		require.Equal(t, data, d)
	})
	t.Run("Has", func(t *testing.T) {
		has, err := db.Has(key)
		require.NoError(t, err)
		require.True(t, has)
	})
	t.Run("Iterator", func(t *testing.T) {
		var d TestDBData
		iter := db.NewIterator()

		require.True(t, iter.Next())
		require.Equal(t, key, iter.Key())
		require.NoError(t, iter.Value(&d))
		require.Equal(t, data, d)

		require.True(t, iter.Next())
		require.Equal(t, key2, iter.Key())
		require.NoError(t, iter.Value(&d))
		require.Equal(t, data2, d)

		require.False(t, iter.Next())
		iter.Release()
		require.NoError(t, iter.Error())
	})
	t.Run("Delete", func(t *testing.T) {
		require.NoError(t, db.Delete(key))
		t.Run("Get", func(t *testing.T) {
			var d TestDBData
			require.Error(t, db.Get(key, &d))
		})
		t.Run("Has", func(t *testing.T) {
			has, err := db.Has(key)
			require.NoError(t, err)
			require.False(t, has)
		})
	})
}

type TestStorePanic struct {
	transient.Store
}

func (s *TestStorePanic) Get(key []byte) []byte {
	panic(errors.New("TestStorePanicGet"))
}
func (s *TestStorePanic) Has(key []byte) bool {
	panic(errors.New("TestStorePanicHas"))
}
func (s *TestStorePanic) Set(key, value []byte) {
	panic(errors.New("TestStorePanicSet"))
}
func (s *TestStorePanic) Delete(key []byte) {
	panic("TestStorePanicDelete")
}

func TestCosmosDbTypedPanic(t *testing.T) {
	cdc := codec.New()
	db := NewDB(&TestStorePanic{*transient.NewStore()}, cdc)
	cdc.RegisterConcrete(&TestDBData{}, "db: TestDBData", nil)
	t.Run("Save", func(t *testing.T) {
		require.EqualError(t, db.Save(nil, TestDBData{A: "test"}), "db: TestStorePanicSet")
	})
	t.Run("Delete", func(t *testing.T) {
		require.EqualError(t, db.Delete(nil), "db: TestStorePanicDelete")
	})
	t.Run("Has", func(t *testing.T) {
		_, err := db.Has(nil)
		require.EqualError(t, err, "db: TestStorePanicHas")
	})
	t.Run("Get", func(t *testing.T) {
		require.EqualError(t, db.Get(nil, nil), "db: TestStorePanicHas")
	})
	t.Run("Iterator", func(t *testing.T) {
		// require.Error(t, db.NewIterator())
		// require.EqualError(t, db.Get(nil, nil), "db: TestStorePanicHas")
	})
}

type TestIteratorPanic struct {
	types.Iterator
}

func (i *TestIteratorPanic) Valid() bool {
	panic(errors.New("TestIteratorPanicValid"))
}
func (i *TestIteratorPanic) Next() {
	panic(errors.New("TestIteratorPanicNext"))
}
func (i *TestIteratorPanic) Key() []byte {
	panic(errors.New("TestIteratorPanicKey"))
}
func (i *TestIteratorPanic) Value() []byte {
	panic(errors.New("TestIteratorPanicValue"))
}
func (i *TestIteratorPanic) Close() {
	panic("TestIteratorPanicClose")
}

func TestDBIteratorPanic(t *testing.T) {
	db := NewDB(transient.NewStore(), nil)
	iter := &DBIterator{
		iter:  &TestIteratorPanic{types.KVStorePrefixIterator(db.store, nil)},
		valid: true,
		cdc:   nil,
	}
	t.Run("Next", func(t *testing.T) {
		require.False(t, iter.Next())
		require.EqualError(t, iter.Error(), "db iterator: TestIteratorPanicNext")
	})
	t.Run("Key", func(t *testing.T) {
		require.Nil(t, iter.Key())
		require.EqualError(t, iter.Error(), "db iterator: TestIteratorPanicKey")
	})
	t.Run("Data", func(t *testing.T) {
		require.NoError(t, iter.Value(nil))
		require.EqualError(t, iter.Error(), "db iterator: TestIteratorPanicValue")
	})
	t.Run("Release", func(t *testing.T) {
		iter.Release()
		require.EqualError(t, iter.Error(), "db iterator: TestIteratorPanicClose")
	})
}
