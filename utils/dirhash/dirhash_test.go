package dirhash

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirHash(t *testing.T) {
	var tests = []struct {
		name  string
		path  string
		extra []byte
		hash  string
	}{
		{
			"regular file",
			"testdata/01",
			nil,
			"0a8a93b56512430570c6790fcf5433c110a52dff92da6cbc6e0a30c474b02941",
		},
		{
			"regular file with read mode permissions",
			"testdata/01_mode",
			nil,
			"bbc8492ba026ee29b1135dc937b65e5be5855108eab6bce8bfcb2d65c832244f",
		},
		{
			"regular file with extra data",
			"testdata/01",
			[]byte{'0'},
			"300a8a93b56512430570c6790fcf5433c110a52dff92da6cbc6e0a30c474b02941",
		},
		{
			"directory with file",
			"testdata/02",
			nil,
			"64cde01ae61c1b66f916f6e2fc7bb1b9fd0e3dbdf5670645ea4d6e3787b9d431",
		},
		{
			"directory with file with extra data",
			"testdata/02",
			[]byte{'0'},
			"3064cde01ae61c1b66f916f6e2fc7bb1b9fd0e3dbdf5670645ea4d6e3787b9d431",
		},
		{
			"symlink",
			"testdata/03",
			nil,
			"c17353643f152f2c43b6810fd45661a83685cfaab0224cb7c37f9fc8a12f409b",
		},
		{
			"symlink with extra data",
			"testdata/03",
			[]byte{'0'},
			"30c17353643f152f2c43b6810fd45661a83685cfaab0224cb7c37f9fc8a12f409b",
		},
		{
			"regular file + directory + symlink",
			"testdata/04",
			nil,
			"d1246aab4e312513c5047abe400f10fc2f2af7d5fd0cf0dc528c0e71732ab913",
		},
		{
			"regular file + directory + symlink with extra data",
			"testdata/04",
			[]byte{'0'},
			"30d1246aab4e312513c5047abe400f10fc2f2af7d5fd0cf0dc528c0e71732ab913",
		},
	}

	for _, tt := range tests {
		hash, err := New(tt.path).Sum(tt.extra)
		assert.NoError(t, err, tt.name)
		assert.Equal(t, tt.hash, hex.EncodeToString(hash), tt.name)
	}
}
