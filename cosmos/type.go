package cosmos

import (
	"fmt"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

// common attribute keys.
const (
	AttributeKeyHash = "hash"
)

// EventHashType is a message with resource hash
var EventHashType = sdktypes.EventTypeMessage + "." + AttributeKeyHash

// EventActionQuery returns tendermint query which matches given message type.
func EventActionQuery(msgType string) string {
	return fmt.Sprintf("%s.%s='%s'", sdktypes.EventTypeMessage, sdktypes.AttributeKeyAction, msgType)
}

// EventModuleQuery returns tendermint query which matches given module.
func EventModuleQuery(module string) string {
	return fmt.Sprintf("%s.%s='%s'", sdktypes.EventTypeMessage, sdktypes.AttributeKeyModule, module)
}
