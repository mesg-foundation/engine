package ethwallet

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/systemservices/ethwallet/x/xgo-ethereum/xaccounts"
)

func (s *Ethwallet) export(inputs *types.Struct) (*types.Struct, error) {
	address := common.HexToAddress(inputs.Fields["address"].GetStringValue())
	account, err := xaccounts.GetAccount(s.keystore, address)
	if err != nil {
		return nil, errAccountNotFound
	}

	passphrase := inputs.Fields["passphrase"].GetStringValue()
	keyJSON, err := s.keystore.Export(account, passphrase, passphrase)
	if err != nil {
		return nil, err
	}

	var accountJSON encryptedKeyJSONV3
	if err = json.Unmarshal(keyJSON, &accountJSON); err != nil {
		return nil, err
	}

	// type CryptoJSON struct {
	// 	Cipher       string                 `json:"cipher"`
	// 	CipherText   string                 `json:"ciphertext"`
	// 	CipherParams cipherparamsJSON       `json:"cipherparams"`
	// 	KDF          string                 `json:"kdf"`
	// 	KDFParams    map[string]interface{} `json:"kdfparams"`
	// 	MAC          string                 `json:"mac"`
	// }

	return &types.Struct{
		Fields: map[string]*types.Value{
			"address": {
				Kind: &types.Value_StringValue{
					StringValue: accountJSON.Address,
				},
			},
			"crypto": {
				Kind: &types.Value_StructValue{
					StructValue: &types.Struct{
						Fields: map[string]*types.Value{
							"cipher": {
								Kind: &types.Value_StringValue{
									StringValue: accountJSON.Crypto.Cipher,
								},
							},
							"ciphertext": {
								Kind: &types.Value_StringValue{
									StringValue: accountJSON.Crypto.CipherText,
								},
							},
							"cipherparams": {
								Kind: &types.Value_StructValue{
									StructValue: &types.Struct{
										Fields: map[string]*types.Value{
											"iv": {
												Kind: &types.Value_StringValue{
													StringValue: accountJSON.Crypto.CipherParams.IV,
												},
											},
										},
									},
								},
							},
							"kdf": {
								Kind: &types.Value_StringValue{
									StringValue: accountJSON.Crypto.KDF,
								},
							},
							"mac": {
								Kind: &types.Value_StringValue{
									StringValue: accountJSON.Crypto.MAC,
								},
							},
						},
					},
				},
			},
			"id": {
				Kind: &types.Value_StringValue{
					StringValue: accountJSON.ID,
				},
			},
			"version": {
				Kind: &types.Value_NumberValue{
					NumberValue: float64(accountJSON.Version),
				},
			},
		},
	}, nil
}
