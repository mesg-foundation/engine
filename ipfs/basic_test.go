package ipfs

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodAndPath(t, r, http.MethodPost, "/add")
		w.Write([]byte("{}"))
	}))
	defer ts.Close()

	_, err := NewClient(ts.URL).Add("", bytes.NewBuffer([]byte{'-'}))
	assert.NoError(t, err)
}

func TestAddBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadAll(r.Body)
		assert.Contains(t, string(b), "a-a")
		w.Write([]byte("{}"))
	}))
	defer ts.Close()

	_, err := NewClient(ts.URL).Add("", bytes.NewBuffer([]byte("a-a")))
	assert.NoError(t, err)
}

func TestGetURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodAndPath(t, r, http.MethodGet, "/get")
	}))
	defer ts.Close()

	_, err := NewClient(ts.URL).Get("")
	assert.NoError(t, err)
}

func TestGetQuery(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "arg=Q", r.URL.RawQuery)
	}))
	defer ts.Close()

	_, err := NewClient(ts.URL).Get("Q")
	assert.NoError(t, err)
}
