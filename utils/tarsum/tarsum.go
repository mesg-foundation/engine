package tarsum

import (
	"archive/tar"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
)

// DefaultHash is default TarSum hashing algorithm - "sha256".
var DefaultHash = sha256.New

// TarSum is the struct for a tar hash calculation.
type TarSum struct {
	stream     *tar.Reader
	fileHashes FileInfoHashes
	newHash    func() hash.Hash
}

// New creates a new struct for calculating a fixed time checksum of a tar archive.
func New(stream *tar.Reader) *TarSum {
	return NewHash(stream, DefaultHash)
}

// NewHash creates a new struct for calculating a fixed time checksum of a tar archive,
// providing a Hash to use rather than the default one.
func NewHash(stream *tar.Reader, newHash func() hash.Hash) *TarSum {
	return &TarSum{
		stream:  stream,
		newHash: newHash,
	}
}

// Sum calculates the sum of tar archive.
func (ts *TarSum) Sum(extra []byte) ([]byte, error) {
	for {
		hdr, err := ts.stream.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("tarsum: %s", err)
		}
		if !isTarFileHashable(hdr.Typeflag) {
			return nil, fmt.Errorf("tarsum: file %s is not hashable", hdr.Name)
		}

		fhash := ts.newHash()
		if _, err := io.Copy(fhash, ts.stream); err != nil {
			return nil, fmt.Errorf("tarsum: %s", err)
		}
		ts.fileHashes = append(ts.fileHashes, FileInfoHash{
			Header: hdr,
			Hash:   fhash.Sum(nil),
		})
	}
	ts.fileHashes.SortByHashes()
	h := ts.newHash()
	for _, fih := range ts.fileHashes {
		h.Write(fih.toBytes())
	}
	return h.Sum(extra), nil
}

// FileHashes returns all the file hashes.
func (ts *TarSum) FileHashes() FileInfoHashes {
	return ts.fileHashes
}

func isTarFileHashable(typeFlag byte) bool {
	return typeFlag == tar.TypeReg ||
		typeFlag == tar.TypeLink ||
		typeFlag == tar.TypeSymlink ||
		typeFlag == tar.TypeDir

}
