package ipfs

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// noop server and handler for testing
var (
	noopHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	noopServer  = httptest.NewServer(noopHandler)

	errorHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	errorServer = httptest.NewServer(errorHandler)
)

func assertMethodAndPath(t *testing.T, r *http.Request, method string, path string) {
	assert.Equal(t, method, r.Method, "method not found")
	assert.Containsf(t, r.URL.Path, path, "invalid url path")
}

func TestResponseStatusNotOk(t *testing.T) {
	_, err := NewClient(errorServer.URL).get(context.Background(), "/", nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "status not ok")
}

func TestDoInvalidMethod(t *testing.T) {
	_, err := NewClient(noopServer.URL).do(context.Background(), "/", "/", nil, nil, nil)
	assert.Error(t, err, "exptected invalid method error")
}

func TestDoInvalidRequest(t *testing.T) {
	_, err := NewClient("").do(context.Background(), "noop", "/", nil, nil, nil)
	assert.Error(t, err, "exptected invalid method error")
}
