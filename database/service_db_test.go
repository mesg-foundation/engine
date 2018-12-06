package database

import (
	"os"
	"sync"
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
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
}

func TestServiceDBSave(t *testing.T) {
	db, closer := openServiceDB(t)
	defer closer()

	s1 := &service.Service{SID: "1", Name: "test-service"}
	require.NoError(t, db.Save(s1))

	// save same service. should replace
	require.NoError(t, db.Save(s1))
	ss, _ := db.All()
	require.Len(t, ss, 1)

	// different sid. should not replace anything.
	s2 := &service.Service{SID: "2", Name: "test-service"}
	require.NoError(t, db.Save(s2))
	ss, _ = db.All()
	require.Len(t, ss, 2)

	// test service without sid.
	s := &service.Service{Name: "test-service"}
	require.EqualError(t, db.Save(s), errCannotSaveWithoutSID.Error())
}

func TestServiceDBGet(t *testing.T) {
	db, closer := openServiceDB(t)
	defer closer()

	want := &service.Service{SID: "2", Name: "test-service"}
	require.NoError(t, db.Save(want))
	defer db.Delete(want.SID)

	// sid.
	got, err := db.Get(want.SID)
	require.NoError(t, err)
	require.Equal(t, want, got)

	// test return err not found.
	_, err = db.Get("3")
	require.Error(t, err)
	require.True(t, IsErrNotFound(err))
}

func TestServiceDBDelete(t *testing.T) {
	db, closer := openServiceDB(t)
	defer closer()

	s := &service.Service{SID: "1", Name: "test-service"}
	require.NoError(t, db.Save(s))
	require.NoError(t, db.Delete(s.SID))
	_, err := db.Get(s.SID)
	require.IsType(t, &ErrNotFound{}, err)
	_, err = db.db.Get([]byte(sidKeyPrefix+s.SID), nil)
	require.Equal(t, leveldb.ErrNotFound, err)
}

func TestServiceDBDeleteConcurrency(t *testing.T) {
	db, closer := openServiceDB(t)
	defer closer()

	s := &service.Service{SID: "2", Name: "test-service"}
	db.Save(s)

	var wg sync.WaitGroup
	errs := make([]error, 0)
	errsM := &sync.Mutex{}
	n := 10
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			if err := db.Delete(s.SID); err != nil {
				errsM.Lock()
				errs = append(errs, err)
				errsM.Unlock()
			}
		}()
	}

	wg.Wait()
	require.Len(t, errs, n-1)
	for i := 0; i < len(errs); i++ {
		require.IsType(t, &ErrNotFound{}, errs[i])
	}
}

func TestServiceDBAll(t *testing.T) {
	db, closer := openServiceDB(t)
	defer closer()

	s1 := &service.Service{SID: "SID1", Name: "test-service"}
	s2 := &service.Service{SID: "SID2", Name: "test-service"}

	require.NoError(t, db.Save(s1))
	require.NoError(t, db.Save(s2))
	defer db.Delete(s1.SID)
	defer db.Delete(s2.SID)

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

	s1 := &service.Service{SID: "2", Name: "test-service"}
	require.NoError(t, db.Save(s1))
	defer db.Delete(s1.SID)

	services, err := db.All()
	require.NoError(t, err)
	require.Len(t, services, 1)
	require.Contains(t, services, s1)
}

func TestIsErrNotFound(t *testing.T) {
	require.True(t, IsErrNotFound(&ErrNotFound{}))
	require.False(t, IsErrNotFound(nil))
}
