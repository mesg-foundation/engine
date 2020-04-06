package orchestrator

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"google.golang.org/grpc/metadata"
)

// RequestSignature is the name of the key to use in the gRPC metadata to set the request signature.
const RequestSignature = "mesg_request_signature"

// Authorizer is the type to aggregate all Admin APIs.
type Authorizer struct {
	cdc               *codec.Codec
	authorizedPubKeys []crypto.PubKey
}

// NewAuthorizer creates a new Authorizer.
func NewAuthorizer(cdc *codec.Codec, authorizedPubKeys []string) (*Authorizer, error) {
	// decode public keys
	pks := make([]crypto.PubKey, 0)
	for _, pkS := range authorizedPubKeys {
		pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, pkS)
		if err != nil {
			return nil, err
		}
		pks = append(pks, pk)
	}

	return &Authorizer{
		cdc:               cdc,
		authorizedPubKeys: pks,
	}, nil
}

// IsAuthorized checks the context for a signature signed by one of the authorizedPubKeys.
func (a *Authorizer) IsAuthorized(ctx context.Context, payload interface{}) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("signature not found in metadata, make sure to set it using the key %q", RequestSignature)
	}
	if len(md[RequestSignature]) == 0 {
		return fmt.Errorf("signature not found in metadata, make sure to set it using the key %q", RequestSignature)
	}
	signature := md[RequestSignature][0]
	signatureBytes, err := base64.RawStdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}
	encodedValue, err := a.cdc.MarshalJSON(payload)
	if err != nil {
		return err
	}
	isAuthorized := false
	for _, authorizedPubKey := range a.authorizedPubKeys {
		if authorizedPubKey.VerifyBytes(encodedValue, signatureBytes) {
			isAuthorized = true
			break
		}
	}
	if !isAuthorized {
		return fmt.Errorf("verification of the signature failed")
	}
	return nil
}
