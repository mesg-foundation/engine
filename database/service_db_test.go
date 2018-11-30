package database

import (
	"os"
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

const testdbname = "db.test"

func openServiceDB(t *testing.T) (*LevelDBServiceDB, func()) {
	deleteDBs(t)
	db, err := NewServiceDB(testdbname)
	require.NoError(t, err)
	return db, func() {
		require.NoError(t, db.Close())
		deleteDBs(t)
	}
}

func deleteDBs(t *testing.T) {
	require.NoError(t, os.RemoveAll(testdbname))
	require.NoError(t, os.RemoveAll(testdbname+aliasesPathSuffix))
}

func TestServiceDBSave(t *testing.T) {
	db, closer := openServiceDB(t)
	defer closer()

	s := &service.Service{ID: "1", Alias: "2", Name: "test-service"}
	require.NoError(t, db.Save(s))
	defer db.Delete(s.ID)

	// test service without id
	s = &service.Service{Name: "test-service", Alias: "alias"}
	require.EqualError(t, db.Save(s), errCannotSaveWithoutID.Error())

	// test service without alias
	s = &service.Service{Name: "test-service", ID: "id"}
	require.EqualError(t, db.Save(s), errCannotSaveWithoutAlias.Error())
}

func TestServiceDBGet(t *testing.T) {
	db, closer := openServiceDB(t)
	defer closer()

	want := &service.Service{ID: "1", Alias: "2", Name: "test-service"}
	require.NoError(t, db.Save(want))
	defer db.Delete(want.ID)

	// id
	got, err := db.Get(want.ID)
	require.NoError(t, err)
	require.Equal(t, want, got)

	// alias
	got, err = db.Get(want.Alias)
	require.NoError(t, err)
	require.Equal(t, want, got)

	// test return err not found
	_, err = db.Get("3")
	require.Error(t, err)
	require.True(t, IsErrNotFound(err))
}

func TestServiceDBDelete(t *testing.T) {
	db, closer := openServiceDB(t)
	defer closer()

	// id
	s := &service.Service{ID: "1", Alias: "2", Name: "test-service"}
	require.NoError(t, db.Save(s))
	require.NoError(t, db.Delete(s.ID))
	_, err := db.Get(s.ID)
	require.True(t, IsErrNotFound(err))

	// alias
	s = &service.Service{ID: "1", Alias: "2", Name: "test-service"}
	require.NoError(t, db.Save(s))
	require.NoError(t, db.Delete(s.Alias))
	_, err = db.Get(s.Alias)
	require.True(t, IsErrNotFound(err))
}

func TestServiceDBAll(t *testing.T) {
	db, closer := openServiceDB(t)
	defer closer()

	s1 := &service.Service{ID: "1", Alias: "alias1", Name: "test-service"}
	s2 := &service.Service{ID: "2", Alias: "alias2", Name: "test-service"}

	require.NoError(t, db.Save(s1))
	require.NoError(t, db.Save(s2))
	defer db.Delete(s1.ID)
	defer db.Delete(s2.ID)

	services, err := db.All()
	require.NoError(t, err)
	require.Len(t, services, 2)
	require.Contains(t, services, s1)
	require.Contains(t, services, s2)
}

func TestServiceDBAllWithDecodeError(t *testing.T) {
	db, closer := openServiceDB(t)
	defer closer()

	id := "idtest"
	require.NoError(t, db.db.Put([]byte(id), []byte("oaiwdhhiodoihwaiohwa"), nil))
	defer db.db.Delete([]byte(id), nil)

	s1 := &service.Service{ID: "1", Alias: "2", Name: "test-service"}
	require.NoError(t, db.Save(s1))
	defer db.Delete(s1.ID)

	services, err := db.All()
	require.NoError(t, err)
	require.Len(t, services, 1)
	require.Contains(t, services, s1)
}

func TestIsErrNotFound(t *testing.T) {
	require.True(t, IsErrNotFound(&ErrNotFound{}))
	require.False(t, IsErrNotFound(nil))
}
