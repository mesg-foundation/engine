package database

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mesg-foundation/core/instance"
	"github.com/stretchr/testify/require"
)

func instancedb(t *testing.T, dir string) InstanceDB {
	db, err := NewInstanceDB(dir)
	require.NoError(t, err)
	return db
}

func TestFindInstance(t *testing.T) {
	dir, _ := ioutil.TempDir("", "TestFindInstance")
	defer os.RemoveAll(dir)
	db := instancedb(t, dir)
	defer db.Close()
	i := &instance.Instance{Hash: "xxx", ServiceHash: "yyy"}
	db.Save(i)
	tests := []struct {
		hash     string
		hasError bool
	}{
		{hash: i.Hash, hasError: false},
		{hash: "yyy", hasError: true},
	}
	for _, test := range tests {
		instance, err := db.Get(test.hash)
		if test.hasError {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
		require.NotNil(t, instance)
		e, err := db.Get(instance.Hash)
		require.NoError(t, err)
		require.NotNil(t, e)
	}
}

func TestSaveInstance(t *testing.T) {
	dir, _ := ioutil.TempDir("", "TestSaveInstance")
	defer os.RemoveAll(dir)
	db := instancedb(t, dir)
	defer db.Close()
	tests := []struct {
		instance *instance.Instance
		hasError bool
	}{
		{&instance.Instance{Hash: "xxx"}, false},
		{&instance.Instance{}, true},
	}
	for _, test := range tests {
		err := db.Save(test.instance)
		if test.hasError {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
	}
}

func TestDeleteInstance(t *testing.T) {
	dir, _ := ioutil.TempDir("", "TestDeleteInstance")
	defer os.RemoveAll(dir)
	db := instancedb(t, dir)
	defer db.Close()
	i := &instance.Instance{Hash: "xxx", ServiceHash: "yyy"}
	db.Save(i)
	require.NoError(t, db.Delete("xxx"))
	inst, err := db.Get("xxx")
	require.Nil(t, inst)
	require.Error(t, err)

	require.NoError(t, db.Delete("yyy"))
}
