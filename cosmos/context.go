package cosmos

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const contextKey = "cosmos-ctx"

func ToContext(ctx sdk.Context) context.Context {
	return context.WithValue(context.Background(), contextKey, ctx)
}

func FromContext(ctx context.Context) (sdk.Context, error) {
	sdkCtx, ok := ctx.Value(contextKey).(sdk.Context)
	if !ok {
		return sdk.Context{}, fmt.Errorf("context is not a cosmos context")
	}
	return sdkCtx, nil
}
