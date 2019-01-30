package ipfs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

const (
	addAPI = "api/v0/add"
)

// IPFS is the struct that contains all informations to access an IPFS node
type IPFS struct {
	client   *http.Client
	endpoint string
}

// Response return the response of a file uploaded
type Response struct {
	Name string
	Hash string
	Size string
}

// New creates a new IPFS client method
func New() *IPFS {
	return &IPFS{
		client:   &http.Client{},
		endpoint: "https://ipfs.infura.io:5001/",
	}
}

// Add data from a reader to the IPFS node
func (ipfs *IPFS) Add(name string, reader io.Reader) (*Response, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("file", name)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, reader); err != nil {
		return nil, err
	}
	w.Close()

	req, err := http.NewRequest("POST", ipfs.endpoint+addAPI, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	res, err := ipfs.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
