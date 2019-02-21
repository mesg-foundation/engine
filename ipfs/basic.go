package ipfs

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

// AddResponse contains the response of adding a file.
type AddResponse struct {
	Name  string
	Hash  string
	Bytes int64
	Size  string
}

// Add sends data from a reader to the ipfs node.
func (ipfs *IPFS) Add(name string, r io.Reader) (*AddResponse, error) {
	b := new(bytes.Buffer)
	w := multipart.NewWriter(b)
	fw, err := w.CreateFormFile("file", name)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, r); err != nil {
		return nil, err
	}
	w.Close()

	headers := http.Header{"Content-Type": []string{w.FormDataContentType()}}
	resp, err := ipfs.post(context.Background(), "add", nil, headers, b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var re AddResponse
	if err := json.NewDecoder(resp.Body).Decode(&re); err != nil {
		return nil, err
	}
	return &re, nil
}

// Get retrives the file from the ipfs node.
func (ipfs *IPFS) Get(hash string) (io.Reader, error) {
	query := url.Values{"arg": []string{hash}}
	resp, err := ipfs.get(context.Background(), "get", query, nil)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
