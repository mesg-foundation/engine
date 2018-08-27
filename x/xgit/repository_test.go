package xgit

import (
	"reflect"
	"testing"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func TestCloneOptions(t *testing.T) {
	for _, tt := range []struct {
		URL string
		o   *git.CloneOptions
	}{
		{"https://localhost", &git.CloneOptions{URL: "https://localhost"}},
		{"localhost", &git.CloneOptions{URL: "https://localhost"}},
		{"localhost#foo", &git.CloneOptions{URL: "https://localhost", ReferenceName: plumbing.ReferenceName("refs/heads/foo")}},
	} {
		options, err := cloneOptions(tt.URL)
		if err != nil {
			t.Errorf("cloneOptions(%q) got error: %s", tt.URL, err)
		}

		if !reflect.DeepEqual(options, tt.o) {
			t.Errorf("cloneOptions(%q) options not equal", tt.URL)
		}
	}

	for _, tt := range []struct {
		URL string
	}{
		{""},
		{"::"},
	} {
		if _, err := cloneOptions(tt.URL); err == nil {
			t.Errorf("cloneOptions(%q) should have failed", tt.URL)
		}
	}
}
