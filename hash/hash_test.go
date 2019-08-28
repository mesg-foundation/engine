package hash

import (
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
}
