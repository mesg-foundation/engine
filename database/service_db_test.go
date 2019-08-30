package database

import (
	"os"
	"sync"
	"testing"

	"github.com/mesg-foundation/engine/database/store"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
	"github.com/stretchr/testify/require"
)

const testdbname = "db.test"

func openServiceDB(t *testing.T) (*ServiceDB, func()) {
	deleteDBs(t)
	store, err := store.NewLevelDBStore(testdbname)
	require.NoError(t, err)
	db := NewServiceDB(store)
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

	s1 := &service.Service{Hash: hash.Int(1)}
	require.NoError(t, db.Save(s1))

	// save same service. should replace
	require.NoError(t, db.Save(s1))
	ss, _ := db.All()
	require.Len(t, ss, 1)

	// different hash, different sid. should not replace anything.
	s3 := &service.Service{Hash: hash.Int(2)}
	require.NoError(t, db.Save(s3))
	ss, _ = db.All()
	require.Len(t, ss, 2)

	// test service without hash.
	require.EqualError(t, db.Save(&service.Service{}), errCannotSaveWithoutHash.Error())
}

func TestServiceDBGet(t *testing.T) {
	db, closer := openServiceDB(t)
	defer closer()

	hs1 := hash.Int(1)

	want := &service.Service{Hash: hs1}
	require.NoError(t, db.Save(want))
	defer db.Delete(hs1)

	// hash.
	got, err := db.Get(hs1)
	require.NoError(t, err)
	require.Equal(t, want, got)

	// test return err not found
	_, err = db.Get(hash.Int(2))
	require.Error(t, err)
}

func TestServiceDBDelete(t *testing.T) {
	db, closer := openServiceDB(t)
	defer closer()

	hs1 := hash.Int(1)

	// hash.
	s := &service.Service{Hash: hs1}
	require.NoError(t, db.Save(s))
	require.NoError(t, db.Delete(hs1))
	_, err := db.Get(hs1)
	require.Error(t, err)
}

// TOFIX: the database is not thread safe anymore...
// Should we lock the db manually? The database could lock the whole db with a mutex.
func TestServiceDBDeleteConcurrency(t *testing.T) {
	t.Skip("delete function need to be fixed or test deleted")
	db, closer := openServiceDB(t)
	defer closer()

	hs1 := hash.Int(1)
	s := &service.Service{Hash: hs1}
	db.Save(s)

	var wg sync.WaitGroup
	errs := make([]error, 0)
	errsM := &sync.Mutex{}
	n := 10
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			if err := db.Delete(hs1); err != nil {
				errsM.Lock()
				errs = append(errs, err)
				errsM.Unlock()
			}
		}()
	}

	wg.Wait()
	require.Len(t, errs, n-1)
	for i := 0; i < len(errs); i++ {
		require.Error(t, errs[i])
	}
}

func TestServiceDBAll(t *testing.T) {
	db, closer := openServiceDB(t)
	defer closer()

	hs1 := hash.Int(1)
	hs2 := hash.Int(2)

	s1 := &service.Service{Hash: hs1}
	s2 := &service.Service{Hash: hs2}

	require.NoError(t, db.Save(s1))
	require.NoError(t, db.Save(s2))
	defer db.Delete(hs1)
	defer db.Delete(hs2)

	services, err := db.All()
	require.NoError(t, err)
	require.Len(t, services, 2)
	require.Contains(t, services, s1)
	require.Contains(t, services, s2)
}

func TestServiceDBAllWithDecodeError(t *testing.T) {
	db, closer := openServiceDB(t)
	defer closer()

	require.NoError(t, db.s.Put(hash.Int(1), []byte("-")))

	services, err := db.All()
	require.NoError(t, err)
	require.Len(t, services, 0)
}
