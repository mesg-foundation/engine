package cosmos

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const requestKey = "cosmos-request"

func ToContext(request sdk.Request) context.Context {
	return context.WithValue(context.Background(), requestKey, request)
}

func ToRequest(ctx context.Context) (sdk.Request, error) {
	request, ok := ctx.Value(requestKey).(sdk.Request)
	if !ok {
		return sdk.Request{}, fmt.Errorf("context is not a cosmos context")
	}
	return request, nil
}
