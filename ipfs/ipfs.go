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

type IPFS struct {
	client   *http.Client
	endpoint string
}

type IPFSResponse struct {
	Name string
	Hash string
	Size string
}

func New() *IPFS {
	return &IPFS{
		client:   &http.Client{},
		endpoint: "https://ipfs.infura.io:5001/",
	}
}

func (ipfs *IPFS) Add(name string, reader io.Reader) (*IPFSResponse, error) {
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
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response IPFSResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
