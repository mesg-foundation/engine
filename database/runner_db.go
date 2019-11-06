package database

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/mesg-foundation/engine/database/store"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/runner"
)

var (
	errCannotSaveRunnerWithoutHash = errors.New("database: can't save runner without hash")
)

// RunnerDB is a database for storing runner definition.
type RunnerDB struct {
	s   store.Store
	cdc *codec.Codec
}

// NewRunnerDB returns the database which is located under given path.
func NewRunnerDB(s store.Store, cdc *codec.Codec) *RunnerDB {
	return &RunnerDB{
		s:   s,
		cdc: cdc,
	}
}

// unmarshal returns the runner from byte slice.
func (d *RunnerDB) unmarshalRunner(hash hash.Hash, value []byte) (*runner.Runner, error) {
	var s runner.Runner
	if err := d.cdc.UnmarshalBinaryBare(value, &s); err != nil {
		return nil, fmt.Errorf("database: could not decode runner %q: %w", hash.String(), err)
	}
	return &s, nil
}

// All returns every runner in database.
func (d *RunnerDB) All() ([]*runner.Runner, error) {
	var (
		runners []*runner.Runner
		iter    = d.s.NewIterator()
	)
	for iter.Next() {
		hash := hash.Hash(iter.Key())
		s, err := d.unmarshalRunner(hash, iter.Value())
		if err != nil {
			return nil, err
		}
		runners = append(runners, s)
	}
	iter.Release()
	return runners, iter.Error()
}

// Save stores runner in database.
// If there is an another runner that uses the same sid, it'll be deleted.
func (d *RunnerDB) Save(r *runner.Runner) error {
	if r.Hash.IsZero() {
		return errCannotSaveRunnerWithoutHash
	}
	b, err := d.cdc.MarshalBinaryBare(r)
	if err != nil {
		return err
	}
	return d.s.Put(r.Hash, b)
}

// Close closes database.
func (d *RunnerDB) Close() error {
	return d.s.Close()
}

// Exist check if runner with given hash exist.
func (d *RunnerDB) Exist(hash hash.Hash) (bool, error) {
	return d.s.Has(hash)
}

// Get retrives runner from database.
func (d *RunnerDB) Get(hash hash.Hash) (*runner.Runner, error) {
	b, err := d.s.Get(hash)
	if err != nil {
		return nil, err
	}
	return d.unmarshalRunner(hash, b)
}

// Delete deletes runner from database.
func (d *RunnerDB) Delete(hash hash.Hash) error {
	return d.s.Delete(hash)
}
