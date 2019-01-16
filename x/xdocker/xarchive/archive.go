package xarchive

import (
	"io"

	"github.com/docker/docker/pkg/archive"
)

// GzippedTar creates an archive from the directory at path compressed with gzip.
func GzippedTar(path string, exclude []string) (io.ReadCloser, error) {
	return archive.TarWithOptions(path, &archive.TarOptions{
		Compression:     archive.Gzip,
		ExcludePatterns: exclude,
	})
}
