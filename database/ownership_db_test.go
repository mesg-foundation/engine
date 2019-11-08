package database

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mesg-foundation/engine/database/store"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/ownership"
	"github.com/stretchr/testify/require"
)

func TestOwnershipDB(t *testing.T) {
	dir, _ := ioutil.TempDir("", "ownership.db.test")
	defer os.RemoveAll(dir)

	store, err := store.NewLevelDBStore(dir)
	require.NoError(t, err)
	db := NewOwnershipDB(store)
	defer db.Close()

	p := &ownership.Ownership{
		Hash:  hash.Int(1),
		Owner: "alice",
		Resource: &ownership.Ownership_ServiceHash{
			ServiceHash: hash.Int(2),
		},
	}

	p2 := &ownership.Ownership{
		Hash:  hash.Int(2),
		Owner: "bob",
		Resource: &ownership.Ownership_ServiceHash{
			ServiceHash: hash.Int(4),
		},
	}

	t.Run("save", func(t *testing.T) {
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
}
