package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDigest(t *testing.T) {
	e := EngineEvent(EngineAPIExecution, map[string]interface{}{
		"foo": "bar",
	})
	assert.Equal(t, e.InstanceHash.String(), "")
	assert.Equal(t, e.Key, "mesg:"+string(EngineAPIExecution))
	assert.Equal(t, e.Data["foo"], "bar")
	assert.Equal(t, e.Hash.String(), "6UW521vQ9KbdDpCYbcKu2FGPLfK8Gsy9kfyF77pzYbYC")
}
