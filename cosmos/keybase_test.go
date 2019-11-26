package cosmos

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeybase(t *testing.T) {
	path, _ := ioutil.TempDir("", "TestKeybase")
	defer os.RemoveAll(path)
	kb, err := NewKeybase(path)
	require.NoError(t, err)

	var (
		name     = "name"
		password = "password"
		mnemonic = ""
	)

	t.Run("Mnemonic", func(t *testing.T) {
		mnemonic, err = kb.NewMnemonic()
		require.NoError(t, err)
	})

	t.Run("Create", func(t *testing.T) {
		acc, err := kb.CreateAccount(name, mnemonic, "", password, 0, 0)
		require.NoError(t, err)
		require.Equal(t, name, acc.GetName())
	})

	t.Run("Exist", func(t *testing.T) {
		exist, err := kb.Exist(name)
		require.NoError(t, err)
		require.True(t, exist)

		exist, err = kb.Exist("random")
		require.NoError(t, err)
		require.False(t, exist)
	})
}
