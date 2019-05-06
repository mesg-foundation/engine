// Package ipfs provides functions to handle ipfs public API.
package ipfs

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/context/ctxhttp"
)

// DefaultEndpoint is endpoint used to access ipfs network.
const DefaultEndpoint = "http://ipfs.app.mesg.com:5001/api/v0"

// DefaultClient is global client used for communication with ipfs network.
var DefaultClient = NewClient(DefaultEndpoint)

// Add sends data from a reader to the ipfs node.
func Add(name string, r io.Reader) (*AddResponse, error) { return DefaultClient.Add(name, r) }

// Get retrives the file from the ipfs node.
func Get(hash string) (io.Reader, error) { return DefaultClient.Get(hash) }

// IPFS handles communication with ipfs network.
type IPFS struct {
	client   *http.Client
	endpoint string
}

// NewClient creates new IPFS client.
func NewClient(endpoint string) *IPFS {
	return &IPFS{
		client:   &http.Client{},
		endpoint: endpoint,
	}
}

// getAPIPath returns the versioned request path to call the api.
// It appends the query parameters to the path if they are not empty.
func (c *IPFS) getAPIPath(path string, query url.Values) string {
	if query == nil {
		return fmt.Sprintf("%s/%s", c.endpoint, path)
	}
	return fmt.Sprintf("%s/%s?%s", c.endpoint, path, query.Encode())
}

func (c *IPFS) get(ctx context.Context, path string, query url.Values, headers http.Header) (*http.Response, error) {
	return c.do(ctx, http.MethodGet, path, query, headers, nil)
}

func (c *IPFS) post(ctx context.Context, path string, query url.Values, headers http.Header, body io.Reader) (*http.Response, error) {
	return c.do(ctx, http.MethodPost, path, query, headers, body)
}

func (c *IPFS) do(ctx context.Context, method, path string, query url.Values, headers http.Header, body io.Reader) (*http.Response, error) {
	fullPath := c.getAPIPath(path, query)
	req, err := http.NewRequest(method, fullPath, body)
	if err != nil {
		return nil, fmt.Errorf("%s %s unknown error: %s", method, fullPath, err)
	}
	for key, value := range headers {
		req.Header[key] = value
	}

	resp, err := ctxhttp.Do(ctx, c.client, req)
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, fmt.Errorf("%s %s i/o timeout", method, fullPath)
		}
		return nil, fmt.Errorf("%s %s unknown error: %s", method, fullPath, err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("%s %s status not ok: %s", method, fullPath, resp.Status)
	}
	return resp, nil
}
