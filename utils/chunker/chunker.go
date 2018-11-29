// Package chunker provides functionalities to create chunks from an
// io.Reader and also create an io.Reader from chunks.
package chunker

import (
	"io"
)

// Chunker converts data read from reader into chunks.
type Chunker struct {
	// chunkSize determines the chunk sizes.
	chunkSize int

	// reader is the reader to convert into chunks.
	reader io.Reader

	// chunks is the channel where upcoming chunks sent.
	chunks chan Data

	// err is the channel where the error will be send after chunking
	// operation fails.
	err chan error

	// value carries the context value of chunk.
	value interface{}

	closed chan struct{}
}

// Data represents a data chunk.
type Data struct {
	// Value carries the context value of data chunk.
	Value interface{}

	// Data is data chunk.
	Data []byte
}

// New returns a new chunker for r and it forwards each chunk to chunks channel.
// an error will be sent to err channel if chunking fails.
func New(r io.Reader, chunks chan Data, err chan error, options ...Option) *Chunker {
	c := &Chunker{
		reader:    r,
		chunks:    chunks,
		err:       err,
		closed:    make(chan struct{}),
		chunkSize: 1024,
	}
	for _, option := range options {
		option(c)
	}
	go c.read()
	return c
}

// Option is the configuration func of Chunker.
type Option func(*Chunker)

// ChunkSizeOption returns an option to set chunk size.
func ChunkSizeOption(n int) Option {
	return func(c *Chunker) {
		c.chunkSize = n
	}
}

// ValueOption returns an option to set context value to chunks.
func ValueOption(value interface{}) Option {
	return func(c *Chunker) {
		c.value = value
	}
}

// read reads data from reader and sends to chunks channel.
func (c *Chunker) read() {
	buf := make([]byte, c.chunkSize)
	for {
		n, err := c.reader.Read(buf)
		if err != nil {
			c.err <- err
			return
		}
		select {
		case <-c.closed:
			return
		case c.chunks <- Data{
			Value: c.value,
			Data:  buf[:n],
		}:
		}
	}
}

// Close will stop reading from reader and no more chunks will be emitted.
func (c *Chunker) Close() error {
	close(c.closed)
	return nil
}
