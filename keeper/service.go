package keeper

// import (
// 	"errors"
// 	"fmt"

// 	"github.com/cosmos/cosmos-sdk/codec"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/mesg-foundation/engine/hash"
// 	"github.com/mesg-foundation/engine/service"
// 	"github.com/sirupsen/logrus"
// 	"github.com/syndtr/goleveldb/leveldb"
// )

// var (
// 	errCannotSaveWithoutHash = errors.New("keeper: can't save service without hash")
// )

// // ServiceKeeper describes the API of database package.
// type ServiceKeeper struct {
// 	storeKey sdk.StoreKey
// 	cdc      *codec.Codec
// }

// // NewServiceKeeper returns the database which is located under given path.
// func NewServiceKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) *ServiceKeeper {
// 	return &ServiceKeeper{
// 		storeKey: storeKey,
// 		cdc:      cdc,
// 	}
// }

// // marshal returns the byte slice from service.
// func (k *ServiceKeeper) marshal(s *service.Service) ([]byte, error) {
// 	return k.cdc.MarshalBinaryBare(s)
// }

// // unmarshal returns the service from byte slice.
// func (k *ServiceKeeper) unmarshal(hash hash.Hash, value []byte) (*service.Service, error) {
// 	var s service.Service
// 	if err := k.cdc.UnmarshalBinaryBare(value, &s); err != nil {
// 		return nil, fmt.Errorf("keeper: could not decode service %q: %s", hash, err)
// 	}
// 	return &s, nil
// }

// // All returns every service in database.
// func (k *ServiceKeeper) All(ctx sdk.Context) ([]*service.Service, error) {
// 	var (
// 		store    = NewStore(ctx.KVStore(k.storeKey))
// 		services []*service.Service
// 		iter     = store.NewIterator()
// 	)
// 	for iter.Next() {
// 		hash := hash.Hash(iter.Key())
// 		s, err := k.unmarshal(hash, iter.Value())
// 		if err != nil {
// 			// NOTE: Ignore all decode errors (possibly due to a service
// 			// structure change or database corruption)
// 			logrus.WithField("service", hash.String()).Warning(err.Error())
// 			continue
// 		}
// 		services = append(services, s)
// 	}
// 	iter.Release()
// 	return services, iter.Error()
// }

// // Delete deletes service from database.
// func (k *ServiceKeeper) Delete(ctx sdk.Context, hash hash.Hash) error {
// 	store := NewStore(ctx.KVStore(k.storeKey))
// 	if _, err := store.Get(hash); err != nil {
// 		// TODO: check not found
// 		// if err == leveldb.ErrNotFound {
// 		// 	return &ErrNotFound{resource: "service", hash: hash}
// 		// }
// 		return err
// 	}
// 	return store.Delete(hash)
// }

// // Get retrives service from database.
// func (k *ServiceKeeper) Get(ctx sdk.Context, hash hash.Hash) (*service.Service, error) {
// 	b, err := d.db.Get(hash, nil)
// 	if err != nil {
// 		if err == leveldb.ErrNotFound {
// 			return nil, &ErrNotFound{resource: "service", hash: hash}
// 		}
// 		return nil, err
// 	}
// 	return k.unmarshal(hash, b)
// }

// // Save stores service in database.
// // If there is an another service that uses the same sid, it'll be deleted.
// func (k *ServiceKeeper) Save(ctx sdk.Context, ctx sdk.Context, s *service.Service) error {
// 	if s.Hash.IsZero() {
// 		return errCannotSaveWithoutHash
// 	}

// 	b, err := d.marshal(s)
// 	if err != nil {
// 		return err
// 	}
// 	return d.db.Put(s.Hash, b, nil)
// }

// // Close closes database.
// func (k *ServiceKeeper) Close() error {
// 	return d.db.Close()
// }

// // ErrNotFound is an not found error.
// type ErrNotFound struct {
// 	hash     hash.Hash
// 	resource string
// }

// func (e *ErrNotFound) Error() string {
// 	return fmt.Sprintf("keeper: %s %q not found", e.resource, e.hash)
// }

// // IsErrNotFound returns true if err is type of ErrNotFound, false otherwise.
// func IsErrNotFound(err error) bool {
// 	_, ok := err.(*ErrNotFound)
// 	return ok
// }
