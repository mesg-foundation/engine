package dirhash

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
)

// DefaultHash is default DirHash hashing algorithm - "sha256".
var DefaultHash = sha256.New

// DirHash is the struct for a tar hash calculation.
type DirHash struct {
	path    string
	newHash func() hash.Hash
}

// New creates a new struct for calculating a fixed time checksum of a tar archive.
func New(path string) *DirHash {
	return NewHash(path, DefaultHash)
}

// NewHash creates a new struct for calculating a fixed time checksum of a tar archive,
// providing a Hash to use rather than the default one.
func NewHash(path string, newHash func() hash.Hash) *DirHash {
	return &DirHash{
		path:    path,
		newHash: newHash,
	}
}

// Hash calculates the hash of directory.
func (ds *DirHash) Sum(extra []byte) ([]byte, error) {
	fhash := ds.newHash()
	if err := filepath.Walk(ds.path, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// make path os-independent to get the same hash on every os
		// despite different path separator.
		safepath := filepath.ToSlash(path)
		fhash.Write([]byte(safepath))

		fm := fi.Mode()

		// write mode in little endian to get constant hash
		// across different cpu architeture.
		binary.Write(fhash, binary.LittleEndian, fm)

		// git dosen't support socket, fifo, block and char devices
		// so only dirs, files and symlinks can be processed.
		switch {
		case fm.IsDir():
			// nothing to do
		case fm.IsRegular():
			// copy file content
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := io.Copy(fhash, f); err != nil {
				return err
			}
		case fm&os.ModeSymlink != 0:
			// copy symlink target
			target, err := os.Readlink(path)
			if err != nil {
				return err
			}
			// see: safepath
			safetarget := filepath.ToSlash(target)
			fhash.Write([]byte(safetarget))
		default:
			return fmt.Errorf("%s invalid file mode", fi.Name())
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("dirhash: %s", err)
	}

	return fhash.Sum(extra), nil
}
