package database

import (
	"os"
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

const testdbname = "db.test"

func openServiceDB(t *testing.T) (*LevelDBServiceDB, func()) {
	db, err := NewServiceDB(testdbname)
	require.NoError(t, err)
	return db, func() {
		require.NoError(t, db.Close())
		require.NoError(t, os.RemoveAll(testdbname))
	}
}

func TestDecodeError(t *testing.T) {
	db, closer := openServiceDB(t)
	defer closer()
	_, err := db.unmarshal("IDToTest", []byte("oaiwdhhiodoihwaiohwa"))
	require.IsType(t, &DecodeError{}, err)
}

func TestServiceDBSave(t *testing.T) {
	db, closer := openServiceDB(t)
	defer closer()

	s := &service.Service{ID: "1", Name: "test-service"}
	require.NoError(t, db.Save(s))
	_, err := db.Get(s.ID)
	require.NoError(t, err)

	// test service without id
	s = &service.Service{Name: "test-service"}
	require.Error(t, db.Save(s))
}

func TestServiceDBGet(t *testing.T) {
	db, closer := openServiceDB(t)
	defer closer()

	want := &service.Service{ID: "1", Name: "test-service"}
	require.NoError(t, db.Save(want))
	got, err := db.Get(want.ID)
	require.NoError(t, err)
	require.Equal(t, want, got)

	// test return err not found
	_, err = db.Get("2")
	require.Error(t, err)
	require.True(t, IsErrNotFound(err))
}

func TestServiceDBDelete(t *testing.T) {
	db, closer := openServiceDB(t)
	defer closer()

	s := &service.Service{ID: "1", Name: "test-service"}
	require.NoError(t, db.Save(s))
	require.NoError(t, db.Delete(s.ID))
	_, err := db.Get(s.ID)
	require.True(t, IsErrNotFound(err))
}

func TestServiceDBAll(t *testing.T) {
	db, closer := openServiceDB(t)
	defer closer()

	s1 := &service.Service{ID: "1", Name: "test-service"}
	s2 := &service.Service{ID: "2", Name: "test-service"}

	require.NoError(t, db.Save(s1))
	require.NoError(t, db.Save(s2))

	services, err := db.All()
	require.NoError(t, err)
	require.Len(t, services, 2)
	require.Contains(t, services, s1)
	require.Contains(t, services, s2)
}

func TestIsErrNotFound(t *testing.T) {
	require.True(t, IsErrNotFound(&ErrNotFound{}))
	require.False(t, IsErrNotFound(nil))
}
