package cosmos

import (
	"fmt"
	"strings"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	eventHashAttr = "hash"
)

func eventHashKey() string {
	return strings.Join([]string{cosmostypes.EventTypeMessage, eventHashAttr}, ".")
}

func EventActionQuery(msgType string) string {
	return fmt.Sprintf("%s.%s='%s'", sdktypes.EventTypeMessage, sdktypes.AttributeKeyAction, msgType)
}

func EventModuleQuery(module string) string {
	return fmt.Sprintf("%s.%s='%s'", sdktypes.EventTypeMessage, sdktypes.AttributeKeyModule, module)
}
