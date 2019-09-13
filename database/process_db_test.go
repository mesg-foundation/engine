package database

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	"github.com/stretchr/testify/require"
)

func TestProcessDB(t *testing.T) {
	dir, _ := ioutil.TempDir("", "process.db.test")
	defer os.Remove(dir)

	db, err := NewProcessDB(dir)
	require.NoError(t, err)
	defer db.Close()

	p := &process.Process{
		Hash: hash.Int(1),
		Key:  "key",
		Edges: []*process.Process_Edge{
			{
				Src: "src",
				Dst: "dst",
			},
		},
		Nodes: []*process.Process_Node{
			{
				Type: &process.Process_Node_Result_{
					Result: &process.Process_Node_Result{
						Key:          "key",
						InstanceHash: hash.Int(1),
						TaskKey:      "taskKey",
					},
				},
			},
			{
				Type: &process.Process_Node_Filter_{
					Filter: &process.Process_Node_Filter{},
				},
			},
		},
	}

	t.Run("save", func(t *testing.T) {
		require.NoError(t, db.Save(p))
	})
	t.Run("get", func(t *testing.T) {
		p1, err := db.Get(p.Hash)
		require.NoError(t, err)
		require.True(t, p.Equal(p1))
	})
	t.Run("all", func(t *testing.T) {
		ps, err := db.All()
		require.NoError(t, err)
		require.Len(t, ps, 1)
		require.True(t, p.Equal(ps[0]))
	})
	t.Run("delete", func(t *testing.T) {
		require.NoError(t, db.Delete(p.Hash))
	})
}

func TestMultipleInsert(t *testing.T) {
	dir, _ := ioutil.TempDir("", "process.db.test")
	defer os.Remove(dir)

	db, err := NewProcessDB(dir)
	require.NoError(t, err)
	defer db.Close()

	p1 := &process.Process{
		Hash: hash.Int(1),
		Key:  "1",
	}

	p2 := &process.Process{
		Hash: hash.Int(2),
		Key:  "2",
	}

	db.Save(p1)
	defer db.Delete(p1.Hash)
	db.Save(p2)
	defer db.Delete(p2.Hash)
	list, err := db.All()
	require.NoError(t, err)
	require.Len(t, list, 2)

	require.Equal(t, list[0].Key, p1.Key)
	require.Equal(t, list[1].Key, p2.Key)

	require.Equal(t, list[0].Hash, p1.Hash)
	require.Equal(t, list[1].Hash, p2.Hash)
}
