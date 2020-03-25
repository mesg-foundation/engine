package main

import (
	"testing"

	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/instance"
	"github.com/stretchr/testify/require"
)

var testInstanceHash hash.Hash

func testInstance(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		var inst *instance.Instance
		lcdGet(t, "instance/get/"+testInstanceHash.String(), &inst)
		require.Equal(t, testInstanceHash, inst.Hash)
		require.Equal(t, testServiceHash, inst.ServiceHash)
		require.Equal(t, hash.Dump([]string{"BAR=3", "FOO=1", "REQUIRED=4"}), inst.EnvHash)
	})

	t.Run("list", func(t *testing.T) {
		insts := make([]*instance.Instance, 0)
		lcdGet(t, "instance/list", &insts)
		require.Len(t, insts, 1)
	})
}
