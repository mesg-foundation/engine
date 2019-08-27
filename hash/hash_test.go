package hash

import (
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	zero       = Int(0)
	one        = Int(1)
	oneStr     = "4uQeVj5tqViQh7yWWGStvkEG1Zmhx6uasJtWCJziofM"
	oneJSONStr = "\"" + oneStr + "\""
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
	assert.Equal(t, uint8(1), Int(1)[0])
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
	hash, err := Decode(oneStr)
	assert.NoError(t, err)
	assert.Equal(t, hash, one)

	_, err = Decode("0")
	assert.Equal(t, "hash: invalid base58 digit ('0')", err.Error())

	_, err = Decode("1")
	assert.Equal(t, "hash: invalid length string", err.Error())
}

func TestIsZero(t *testing.T) {
	assert.True(t, Hash{}.IsZero())
}

func TestString(t *testing.T) {
	assert.Equal(t, one.String(), oneStr)
}

func TestEqual(t *testing.T) {
	assert.True(t, zero.Equal(zero))
	assert.False(t, zero.Equal(one))
}

func TestMarshal(t *testing.T) {
	b, err := one.Marshal()
	assert.NoError(t, err)
	assert.Equal(t, oneStr, string(b))
}

func TestUnmarshal(t *testing.T) {
	var h Hash
	assert.NoError(t, h.Unmarshal([]byte(oneStr)))
	assert.Equal(t, one, h)
}

func TestSize(t *testing.T) {
	assert.Equal(t, 0, Hash(nil).Size())
	assert.Equal(t, size, zero.Size())
}

func TestMarshalJSON(t *testing.T) {
	b, err := one.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, oneJSONStr, string(b))
}

func TestUnmarshalJSON(t *testing.T) {
	var h Hash
	assert.NoError(t, h.UnmarshalJSON([]byte(oneJSONStr)))
	assert.Equal(t, one, h)
}
