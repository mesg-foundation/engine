package xgit

import (
	"errors"
	"net/url"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// Clone clones a repository from url into the path.
func Clone(URL string, path string) error {
	options, err := cloneOptions(URL)
	if err != nil {
		return err
	}
	_, err = git.PlainClone(path, false, options)
	return err
}

func cloneOptions(URL string) (*git.CloneOptions, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	if u.String() == "" {
		return nil, errors.New("empty url")
	}

	if u.Scheme == "" {
		u.Scheme = "https"
	}

	options := &git.CloneOptions{}
	if u.Fragment != "" {
		options.ReferenceName = plumbing.ReferenceName("refs/heads/" + u.Fragment)
		u.Fragment = ""
	}
	options.URL = u.String()
	return options, err
}
