package hash

import (
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDigest(t *testing.T) {
	d := New()
	d.Write([]byte{0})

	hash := d.Sum(nil)
	assert.Equal(t, hash.String(), "8RBsoeyoRwajj86MZfZE6gMDJQVYGYcdSfx1zxqxNHbr")

	_, err := Decode(hash.String())
	assert.NoError(t, err)
}

func TestDump(t *testing.T) {
	assert.Equal(t, Dump(struct{}{}).String(), "5ajuwjHoLj33yG5t5UFsJtUb3vnRaJQEMPqSLz6VyoHK")
}

func TestInt(t *testing.T) {
	assert.Len(t, Int(1), size)
}

func TestRandom(t *testing.T) {
	hash, _ := Random()
	assert.Len(t, hash, size)

	seen := make(map[string]bool)

	f := func() bool {
		hash, err := Random()
		if err != nil {
			return false
		}

		if seen[hash.String()] {
			return false
		}
		seen[hash.String()] = true
		return true
	}
	require.NoError(t, quick.Check(f, nil))
}

func TestDecode(t *testing.T) {
	hash, err := Decode("4uQeVj5tqViQh7yWWGStvkEG1Zmhx6uasJtWCJziofM")
	assert.NoError(t, err)
	assert.Equal(t, hash, Int(1))

	_, err = Decode("0")
	assert.Equal(t, "hash: invalid base58 digit ('0')", err.Error())

	_, err = Decode("1")
	assert.Equal(t, "hash: invalid length string", err.Error())
}

func TestIsZero(t *testing.T) {
	assert.True(t, Hash{}.IsZero())
}

func TestString(t *testing.T) {
	assert.Equal(t, Int(1).String(), "4uQeVj5tqViQh7yWWGStvkEG1Zmhx6uasJtWCJziofM")
}

func TestEqual(t *testing.T) {
	assert.True(t, Int(0).Equal(Int(0)))
	assert.False(t, Int(0).Equal(Int(1)))
}
