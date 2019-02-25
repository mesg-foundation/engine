package tarsum

import (
	"archive/tar"
	"bytes"
	"sort"
)

// FileInfoHash provides a struct for accessing file hash
// information within a tar file.
type FileInfoHash struct {
	Header *tar.Header
	Hash   []byte
}

// toBytes serialize FileInfoHash to byte slice for hash calculations.
func (fih *FileInfoHash) toBytes() []byte {
	buf := &bytes.Buffer{}
	buf.Write(fih.Hash)
	buf.Write([]byte(fih.Header.Name))
	buf.Write([]byte{fih.Header.Typeflag})
	if fih.Header.Typeflag == tar.TypeLink ||
		fih.Header.Typeflag == tar.TypeSymlink {
		buf.Write([]byte(fih.Header.Linkname))
	}
	return buf.Bytes()
}

// FileInfoHashes provides a list of FileInfoHash.
type FileInfoHashes []FileInfoHash

// Len returns the size of the FileInfoHashes.
func (fih FileInfoHashes) Len() int { return len(fih) }

// Swap swaps two FileInfoHash values.
func (fih FileInfoHashes) Swap(i, j int) { fih[i], fih[j] = fih[j], fih[i] }

// SortByNames sorts FileInfoHashes content by name.
func (fih FileInfoHashes) SortByNames() {
	sort.Sort(byName{fih})
}

// SortByHashes sorts FileInfoHashes content by sums.
func (fih FileInfoHashes) SortByHashes() {
	sort.Sort(byHash{fih})
}

// byName is a sort.Sort helper for sorting by file names.
// If names are the same, order them by their appearance in the tar archive
type byName struct{ FileInfoHashes }

func (bn byName) Less(i, j int) bool {
	return bn.FileInfoHashes[i].Header.Name < bn.FileInfoHashes[j].Header.Name
}

// byHash is a sort.Sort helper for sorting by the sums of all the fileinfos in the tar archive
type byHash struct{ FileInfoHashes }

func (bs byHash) Less(i, j int) bool {
	return bytes.Compare(bs.FileInfoHashes[i].Hash, bs.FileInfoHashes[j].Hash) < 0
}
