package ethwallet

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/systemservices/ethwallet/x/xgo-ethereum/xaccounts"
)

func (s *Ethwallet) export(inputs *types.Struct) (*types.Struct, error) {
	address := common.HexToAddress(inputs.GetStringValue("address"))
	account, err := xaccounts.GetAccount(s.keystore, address)
	if err != nil {
		return nil, errAccountNotFound
	}

	passphrase := inputs.GetStringValue("passphrase")
	keyJSON, err := s.keystore.Export(account, passphrase, passphrase)
	if err != nil {
		return nil, err
	}

	var accountJSON encryptedKeyJSONV3
	if err = json.Unmarshal(keyJSON, &accountJSON); err != nil {
		return nil, err
	}

	return types.NewStruct(map[string]*types.Value{
		"address": types.NewValueFrom(accountJSON.Address),
		"version": types.NewValueFrom(accountJSON.Version),
		"id":      types.NewValueFrom(accountJSON.ID),
		"crypto": types.NewValueFrom(types.NewStruct(
			map[string]*types.Value{
				"cipher":     types.NewValueFrom(accountJSON.Crypto.Cipher),
				"ciphertext": types.NewValueFrom(accountJSON.Crypto.CipherText),
				"cipherparams": types.NewValueFrom(types.NewStruct(
					map[string]*types.Value{
						"iv": types.NewValueFrom(accountJSON.Crypto.CipherParams.IV),
					},
				)),
				"kdfparams": types.NewValueFrom(accountJSON.Crypto.KDFParams),
				"kdf":       types.NewValueFrom(accountJSON.Crypto.KDF),
				"mac":       types.NewValueFrom(accountJSON.Crypto.MAC),
			},
		)),
	}), nil
}
