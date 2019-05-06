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
			"11d376d7bfad120831035cda21097f8d4a3a3b8cac03c5dc52863a574b1d5614",
		},
		{
			"regular file with read mode permissions",
			"testdata/01_mode",
			nil,
			"7dfd06256cd04975545aa2488af3aac91b6e44b32c5f24fd9f0d508f9a194aa4",
		},
		{
			"regular file with extra data",
			"testdata/01",
			[]byte{'0'},
			"4a8fba57a6d0b0609cd9b0219fca9afecf88ec827e888c772eeda23989a5b0aa",
		},
		{
			"directory with file",
			"testdata/02",
			nil,
			"0a8fc06675a7ee082161b9c85d658b58952c3d22668fc6e84146efc4aef56b3c",
		},
		{
			"directory with file with extra data",
			"testdata/02",
			[]byte{'0'},
			"12cf85424d527741d530632687ae38c63ca5bd0adb7803fe09c7c1a87f21b270",
		},
		{
			"symlink",
			"testdata/03",
			nil,
			"599b534f45d7d986aede06bc12e23ebbd2fe917e3b66f9f4d5a49e646ef42005",
		},
		{
			"symlink with extra data",
			"testdata/03",
			[]byte{'0'},
			"c45664de7a405ef1010ea74066c0572efcf3fa8137ed10d96d0c5aaf8dcd9a63",
		},
		{
			"regular file + directory + symlink",
			"testdata/04",
			nil,
			"b4388f782e77ebfb584bac2ba13396bbb21138c954ad1836163b03bca182099d",
		},
		{
			"regular file + directory + symlink with extra data",
			"testdata/04",
			[]byte{'0'},
			"203fed7ec98926b4926cc6977b0ba6d5e9a3a9325ad5bfa0edcbce671243325c",
		},
	}

	for _, tt := range tests {
		hash, err := New(tt.path).Sum(tt.extra)
		assert.NoError(t, err, tt.name)
		assert.Equal(t, 32, len(hash))
		assert.Equal(t, tt.hash, hex.EncodeToString(hash), tt.name)
	}
}
