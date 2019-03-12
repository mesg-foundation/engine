package service

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/mesg-foundation/core/x/xerrors"
	"github.com/stretchr/testify/require"
)

func TestReadDefinition(t *testing.T) {
	_, err := ReadDefinition("testdata")
	require.NoError(t, err)
}

func TestDecodeStrictDefinition(t *testing.T) {
	buf := bytes.NewBufferString("unknown-filed:")
	_, err := DecodeDefinition(buf)
	require.Error(t, err)
}

func TestValidateDefinition(t *testing.T) {
	mesgfile := filepath.Join("testdata", "mesg.invalid.yml")
	f, err := os.Open(mesgfile)
	require.NoError(t, err)

	service, err := DecodeDefinition(f)
	require.NoError(t, err)

	err = ValidateDefinition(service)
	require.Len(t, err.(xerrors.Errors), 22)
}
