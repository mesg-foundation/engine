package definition

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/mesg-foundation/core/x/xerrors"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	_, err := Read("testdata")
	require.NoError(t, err)
}

func TestDecodeStrict(t *testing.T) {
	buf := bytes.NewBufferString("unknown-filed:")
	_, err := Decode(buf)
	require.Error(t, err)
}

func TestValidate(t *testing.T) {
	mesgfile := filepath.Join("testdata", "mesg.invalid.yml")
	f, err := os.Open(mesgfile)
	require.NoError(t, err)

	service, err := Decode(f)
	require.NoError(t, err)

	err = Validate(service)
	require.Len(t, err.(xerrors.Errors), 25)
}
