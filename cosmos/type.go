package cosmos

import (
	"fmt"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

// common attribute keys.
const (
	AttributeKeyHash = "hash"
)

// EventHashType
var EventHashType = cosmostypes.EventTypeMessage + "." + AttributeKeyHash

func EventActionQuery(msgType string) string {
	return fmt.Sprintf("%s.%s='%s'", sdktypes.EventTypeMessage, sdktypes.AttributeKeyAction, msgType)
}

func EventModuleQuery(module string) string {
	return fmt.Sprintf("%s.%s='%s'", sdktypes.EventTypeMessage, sdktypes.AttributeKeyModule, module)
}
