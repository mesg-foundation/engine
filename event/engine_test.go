package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDigest(t *testing.T) {
	e := EngineEvent(EngineAPIExecution, map[string]interface{}{
		"foo": "bar",
	})
	assert.Equal(t, e.InstanceHash.String(), engineEventHash)
	assert.Equal(t, e.Key, string(EngineAPIExecution))
	assert.Equal(t, e.Data["foo"], "bar")
	assert.Equal(t, e.Hash.String(), "5oHsDq9Njo4uRmjo6ehD3FREYuH6o6PbELn6vYahHAay")
}
