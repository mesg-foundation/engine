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

	s1 := &service.Service{Hash: "00", SID: "1", Name: "test-service"}
	require.NoError(t, db.Save(s1))

	// save same service. should replace
	require.NoError(t, db.Save(s1))
	ss, _ := db.All()
	require.Len(t, ss, 1)

	// different id, same SID. should replace s1
	s2 := &service.Service{Hash: "01", SID: "1", Name: "test-service"}
	require.NoError(t, db.Save(s2))
	_, err := db.Get(s1.Hash)
	require.IsType(t, &ErrNotFound{}, err)

	// different id, different SID. should not replace anything
	s3 := &service.Service{Hash: "02", SID: "2", Name: "test-service"}
	require.NoError(t, db.Save(s3))
	ss, _ = db.All()
	require.Len(t, ss, 2)

	// test service without id
	s := &service.Service{Name: "test-service", SID: "SID"}
	require.EqualError(t, db.Save(s), errCannotSaveWithoutHash.Error())

	// test service without SID
	s = &service.Service{Name: "test-service", Hash: "id"}
	require.EqualError(t, db.Save(s), errCannotSaveWithoutSID.Error())

	// test service id same length as SID
	s = &service.Service{Name: "test-service", Hash: "sameLength", SID: "sameLength"}
	require.EqualError(t, db.Save(s), errSIDSameLen.Error())
}

func TestServiceDBGet(t *testing.T) {
	db, closer := openServiceDB(t)
	defer closer()

	want := &service.Service{Hash: "00", SID: "2", Name: "test-service"}
	require.NoError(t, db.Save(want))
	defer db.Delete(want.Hash)

	// id
	got, err := db.Get(want.Hash)
	require.NoError(t, err)
	require.Equal(t, want, got)

	// SID
	got, err = db.Get(want.SID)
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
	s := &service.Service{Hash: "00", SID: "2", Name: "test-service"}
	require.NoError(t, db.Save(s))
	require.NoError(t, db.Delete(s.Hash))
	_, err := db.Get(s.Hash)
	require.IsType(t, &ErrNotFound{}, err)
	_, err = db.Get(s.SID)
	require.IsType(t, &ErrNotFound{}, err)

	_, err = db.db.Get([]byte(hashKeyPrefix+s.Hash), nil)
	require.Equal(t, leveldb.ErrNotFound, err)
	_, err = db.db.Get([]byte(sidKeyPrefix+s.SID), nil)
	require.Equal(t, leveldb.ErrNotFound, err)

	// SID
	s = &service.Service{Hash: "11", SID: "3", Name: "test-service"}
	require.NoError(t, db.Save(s))
	require.NoError(t, db.Delete(s.SID))
	_, err = db.Get(s.SID)
	require.IsType(t, &ErrNotFound{}, err)
	_, err = db.Get(s.Hash)
	require.IsType(t, &ErrNotFound{}, err)

	_, err = db.db.Get([]byte(hashKeyPrefix+s.Hash), nil)
	require.Equal(t, leveldb.ErrNotFound, err)
	_, err = db.db.Get([]byte(sidKeyPrefix+s.SID), nil)
	require.Equal(t, leveldb.ErrNotFound, err)
}

func TestServiceDBDeleteConcurrency(t *testing.T) {
	db, closer := openServiceDB(t)
	defer closer()

	s := &service.Service{Hash: "00", SID: "2", Name: "test-service"}
	db.Save(s)

	var wg sync.WaitGroup
	errs := make([]error, 0)
	errsM := &sync.Mutex{}
	n := 10
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			if err := db.Delete(s.Hash); err != nil {
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

	s1 := &service.Service{Hash: "00", SID: "SID1", Name: "test-service"}
	s2 := &service.Service{Hash: "01", SID: "SID2", Name: "test-service"}

	require.NoError(t, db.Save(s1))
	require.NoError(t, db.Save(s2))
	defer db.Delete(s1.Hash)
	defer db.Delete(s2.Hash)

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

	s1 := &service.Service{Hash: "00", SID: "2", Name: "test-service"}
	require.NoError(t, db.Save(s1))
	defer db.Delete(s1.Hash)

	services, err := db.All()
	require.NoError(t, err)
	require.Len(t, services, 1)
	require.Contains(t, services, s1)
}

func TestIsErrNotFound(t *testing.T) {
	require.True(t, IsErrNotFound(&ErrNotFound{}))
	require.False(t, IsErrNotFound(nil))
}
