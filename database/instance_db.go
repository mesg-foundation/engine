package database

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/mesg-foundation/engine/database/store"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/instance"
)

var (
	errCannotSaveInstanceWithoutHash = errors.New("database: can't save instance without hash")
)

// InstanceDB is a database for storing instance definition.
type InstanceDB struct {
	s   store.Store
	cdc *codec.Codec
}

// NewInstanceDB returns the database which is located under given path.
func NewInstanceDB(s store.Store, cdc *codec.Codec) *InstanceDB {
	return &InstanceDB{
		s:   s,
		cdc: cdc,
	}
}

// unmarshal returns the instance from byte slice.
func (d *InstanceDB) unmarshalInstance(hash hash.Hash, value []byte) (*instance.Instance, error) {
	var s instance.Instance
	if err := d.cdc.UnmarshalBinaryBare(value, &s); err != nil {
		return nil, fmt.Errorf("database: could not decode instance %q: %w", hash.String(), err)
	}
	return &s, nil
}

// All returns every instance in database.
func (d *InstanceDB) All() ([]*instance.Instance, error) {
	var (
		instances []*instance.Instance
		iter      = d.s.NewIterator()
	)
	for iter.Next() {
		hash := hash.Hash(iter.Key())
		s, err := d.unmarshalInstance(hash, iter.Value())
		if err != nil {
			return nil, err
		}
		instances = append(instances, s)
	}
	iter.Release()
	return instances, iter.Error()
}

// Save stores instance in database.
// If there is an another instance that uses the same sid, it'll be deleted.
func (d *InstanceDB) Save(r *instance.Instance) error {
	if r.Hash.IsZero() {
		return errCannotSaveInstanceWithoutHash
	}
	b, err := d.cdc.MarshalBinaryBare(r)
	if err != nil {
		return err
	}
	return d.s.Put(r.Hash, b)
}

// Close closes database.
func (d *InstanceDB) Close() error {
	return d.s.Close()
}

// Exist check if instance with given hash exist.
func (d *InstanceDB) Exist(hash hash.Hash) (bool, error) {
	return d.s.Has(hash)
}

// Get retrives instance from database.
func (d *InstanceDB) Get(hash hash.Hash) (*instance.Instance, error) {
	b, err := d.s.Get(hash)
	if err != nil {
		return nil, err
	}
	return d.unmarshalInstance(hash, b)
}

// Delete deletes instance from database.
func (d *InstanceDB) Delete(hash hash.Hash) error {
	return d.s.Delete(hash)
}
