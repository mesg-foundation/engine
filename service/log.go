package service

import (
	"io"
)

// LogReader is a struct that holds dependency
// key and logs reader.
type LogReader struct {
	key string
	r   io.ReadCloser
}

// Dependency returns dependency key.
func (r *LogReader) Dependency() string {
	return r.key
}

// Reader returns underlying reader.
func (r *LogReader) Reader() io.ReadCloser {
	return r.r
}

// Close closes underlying reader.
func (r *LogReader) Close() error {
	return r.r.Close()
}
