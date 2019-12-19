package cosmos

import (
	"sync"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/x/auth"
	sdkauth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/mesg-foundation/engine/x/xmath"
)

const (
	// AccNumber is the account number of the account in keybase.
	AccNumber = 0
	// AccIndex is the account index of the account in keybase.
	AccIndex = 0
)

type Account struct {
	auth.Account

	querier sdkauth.NodeQuerier
	mx      sync.Mutex
}

func NewAccount(acc auth.Account, querier sdkauth.NodeQuerier) *Account {
	return &Account{Account: acc, querier: querier}
}

func NewAccountFromKeybase(kb keys.Keybase, name string, querier sdkauth.NodeQuerier) (*Account, error) {
	info, err := kb.Get(name)
	if err != nil {
		return nil, err
	}
	acc := auth.NewBaseAccount(
		info.GetAddress(),
		nil,
		info.GetPubKey(),
		AccNumber,
		0,
	)

	return &Account{Account: acc, querier: querier}, nil
}

func (a *Account) GetSequence() uint64 {
	a.mx.Lock()
	seq := a.Account.GetSequence()
	if acc, err := auth.NewAccountRetriever(a.querier).GetAccount(a.GetAddress()); err == nil {
		seq = xmath.MaxUint64(seq, acc.GetSequence())
		a.Account = acc
	}
	a.SetSequence(seq + 1)
	a.mx.Unlock()
	return seq
}
