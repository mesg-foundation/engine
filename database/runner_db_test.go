package database

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/mesg-foundation/engine/database/store"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/runner"
	"github.com/stretchr/testify/require"
)

func TestRunnerDB(t *testing.T) {
	cdc := codec.New()

	dir, _ := ioutil.TempDir("", "runner.db.test")
	defer os.RemoveAll(dir)

	store, err := store.NewLevelDBStore(dir)
	require.NoError(t, err)
	db := NewRunnerDB(store, cdc)
	defer db.Close()

	p := &runner.Runner{
		Hash:         hash.Int(1),
		Address:      "alice",
		InstanceHash: hash.Int(11),
	}

	p2 := &runner.Runner{
		Hash:         hash.Int(2),
		Address:      "bob",
		InstanceHash: hash.Int(22),
	}

	t.Run("save", func(t *testing.T) {
		require.Error(t, db.Save(&runner.Runner{}))
		require.NoError(t, db.Save(p))
		require.NoError(t, db.Save(p2))
	})
	t.Run("all", func(t *testing.T) {
		ps, err := db.All()
		require.NoError(t, err)
		require.Len(t, ps, 2)
		require.True(t, p.Equal(ps[0]))
		require.True(t, p2.Equal(ps[1]))
	})
	t.Run("exist", func(t *testing.T) {
		exist, err := db.Exist(hash.Int(1))
		require.NoError(t, err)
		require.True(t, exist)
	})
	t.Run("get", func(t *testing.T) {
		get, err := db.Get(hash.Int(1))
		require.NoError(t, err)
		require.True(t, p.Equal(get))
	})
	t.Run("delete", func(t *testing.T) {
		require.NoError(t, db.Delete(hash.Int(1)))
		t.Run("does not exist", func(t *testing.T) {
			exist, err := db.Exist(hash.Int(1))
			require.NoError(t, err)
			require.False(t, exist)
		})
		t.Run("get not existing", func(t *testing.T) {
			_, err := db.Get(hash.Int(1))
			require.Error(t, err)
		})
		t.Run("all after delete", func(t *testing.T) {
			ps, err := db.All()
			require.NoError(t, err)
			require.Len(t, ps, 1)
			require.True(t, p2.Equal(ps[0]))
		})
	})
}
