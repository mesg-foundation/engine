package xrequire

import (
	"io"

	"github.com/docker/docker/pkg/archive"
)

// GzipCompress creates an gzip archive from the directory at path.
func GzipCompress(path string) (io.Reader, error) {
	return archive.TarWithOptions(path, &archive.TarOptions{
		Compression: archive.Gzip,
	})
}
