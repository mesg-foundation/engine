package tarsum

import (
	"archive/tar"
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTarSum(t *testing.T) {
	type tarfile struct {
		header *tar.Header
		body   string
	}

	tests := []struct {
		name  string
		files []tarfile
		extra []byte
		hash  string
	}{
		{
			name: "no file",
			hash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name: "empty file",
			files: []tarfile{
				{
					header: &tar.Header{
						Typeflag: tar.TypeReg,
						Name:     "a",
					},
				},
			},
			hash: "d1edffdff77dfee899d94f5736a5c067f6a9671cc099910ddd4f9490136e9d25",
		},
		{
			name: "empty dir",
			files: []tarfile{
				{
					header: &tar.Header{
						Typeflag: tar.TypeDir,
						Name:     "a",
					},
				},
			},
			hash: "4c6be43c4e542522c4aa264ad7a562f027b7bc4b47543dd691b7782edb1dea36",
		},
		{
			name: "empty link",
			files: []tarfile{
				{
					header: &tar.Header{
						Typeflag: tar.TypeLink,
						Name:     "a",
					},
				},
			},
			hash: "609946c7be8104ffd7b9c5eaf5c7c38ee96f0bb2098cb05276d310107d9cf273",
		},
		{
			name: "empty symlink",
			files: []tarfile{
				{
					header: &tar.Header{
						Typeflag: tar.TypeSymlink,
						Name:     "a",
					},
				},
			},
			hash: "36f9fa5256a99d503fa19490ef1b4d65e70d515f18e2972cd96fa273522002e3",
		},
		{
			name: "file",
			files: []tarfile{
				{
					header: &tar.Header{
						Typeflag: tar.TypeReg,
						Name:     "a",
						Size:     1,
					},
					body: "a",
				},
			},
			hash: "4649927877fb7028b94c4ff3158509ec46e0b7910c19702fe860677c094dd14f",
		},
	}

	for _, tt := range tests {
		buf := &bytes.Buffer{}
		tw := tar.NewWriter(buf)
		for _, file := range tt.files {
			require.NoError(t, tw.WriteHeader(file.header), tt.name)
			_, err := tw.Write([]byte(file.body))
			require.NoError(t, err, tt.name)
		}
		require.NoError(t, tw.Close(), tt.name)

		hash, err := New(tar.NewReader(buf)).Sum(tt.extra)
		require.NoError(t, err, tt.name)
		assert.Equal(t, tt.hash, hex.EncodeToString(hash), tt.name)
	}
}
