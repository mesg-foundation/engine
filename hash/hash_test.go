package hash

import (
	"encoding/json"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	zero = Int(0)
	one  = Int(1)
)

func TestDigest(t *testing.T) {
	d := New()
	d.Write([]byte{0})

	hash := d.Sum(nil)
	assert.Equal(t, hash.String(), "8RBsoeyoRwajj86MZfZE6gMDJQVYGYcdSfx1zxqxNHbr")
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
	hash, err := Decode("4uQeVj5tqViQh7yWWGStvkEG1Zmhx6uasJtWCJziofM")
	assert.NoError(t, err)
	assert.Equal(t, hash, one)

	_, err = Decode("0")
	assert.Equal(t, "hash: invalid base58 digit ('0')", err.Error())

	_, err = Decode("1")
	assert.Equal(t, "hash: invalid length", err.Error())
}

func TestIsZero(t *testing.T) {
	assert.True(t, Hash{}.IsZero())
}

func TestString(t *testing.T) {
	assert.Equal(t, one.String(), "4uQeVj5tqViQh7yWWGStvkEG1Zmhx6uasJtWCJziofM")
}

func TestEqual(t *testing.T) {
	assert.True(t, zero.Equal(zero))
	assert.False(t, zero.Equal(one))
}

func TestSize(t *testing.T) {
	assert.Equal(t, 0, Hash(nil).Size())
	assert.Equal(t, size, zero.Size())
	assert.Equal(t, 5, Hash([]byte("hello")).Size())
}

func TestMarshalJSON(t *testing.T) {
	b, err := json.Marshal(Int(1))
	assert.NoError(t, err)
	assert.Equal(t, "\"4uQeVj5tqViQh7yWWGStvkEG1Zmhx6uasJtWCJziofM\"", string(b))
}

func TestUnmarshalJSON(t *testing.T) {
	var h Hash
	assert.NoError(t, json.Unmarshal([]byte("\"4uQeVj5tqViQh7yWWGStvkEG1Zmhx6uasJtWCJziofM\""), &h))
	assert.Equal(t, Int(1), h)
}

func TestUnmarshal(t *testing.T) {
	var hash Hash
	data := []byte(Int(1))
	hash.Unmarshal(data)
	// check if unmarshal copy the data
	// test if two slises do not share the same address
	assert.True(t, &hash[cap(hash)-1] != &data[cap(data)-1])
}

func TestWrongLength(t *testing.T) {
	var h Hash
	wrongLenByte := []byte("hello")
	t.Run("Marshal", func(t *testing.T) {
		_, err := Hash(wrongLenByte).Marshal()
		require.EqualError(t, err, "hash: invalid length")
	})
	t.Run("MarshalTo", func(t *testing.T) {
		_, err := h.MarshalTo(wrongLenByte)
		require.EqualError(t, err, "hash: invalid length")
	})
	t.Run("Unmarshal", func(t *testing.T) {
		require.EqualError(t, h.Unmarshal(wrongLenByte), "hash: invalid length")
	})
}
