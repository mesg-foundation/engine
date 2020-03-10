package cosmos

import (
	"crypto/sha256"
	"io/ioutil"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
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
		acc, err := kb.CreateAccount(name, mnemonic, "", password, keys.CreateHDPath(0, 0).String(), DefaultAlgo)
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

	t.Run("Sign", func(t *testing.T) {
		name2 := "name2"
		mnemonic2, _ := kb.NewMnemonic()
		_, err := kb.CreateAccount(name2, mnemonic2, "", password, keys.CreateHDPath(0, 0).String(), DefaultAlgo)
		require.NoError(t, err)
		hash := sha256.Sum256([]byte(name + ":" + password))
		hash2 := sha256.Sum256([]byte(name2 + ":" + password))
		for i := 0; i < 1000; i++ {
			nameToUse := name
			if i%2 == 0 {
				nameToUse = name2
			}
			sig, pub, err := kb.Sign(nameToUse, password, []byte{})
			require.NoError(t, err)
			require.NotEmpty(t, sig)
			require.NotEmpty(t, pub)
		}
		require.NotEmpty(t, kb.privKeysCache[hash])
		require.NotEmpty(t, kb.privKeysCache[hash2])
		require.NotEqual(t, kb.privKeysCache[hash], kb.privKeysCache[hash2])
	})
}
