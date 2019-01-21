package xgit

import (
	"errors"
	"net/url"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/protocol/packp"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/client"
)

const defaultScheme = "https"

// IsGitURL checks if given URL is git repository.
func IsGitURL(URL string) bool {
	if strings.HasPrefix(URL, "git://") || strings.HasSuffix(URL, ".git") {
		return true
	}
	_, err := ListRemote(URL)
	return err == nil
}

// ListRemote lists all remotes from given url.
func ListRemote(URL string) (*packp.AdvRefs, error) {
	URL = addDefaultURLScheme(URL)
	ep, err := transport.NewEndpoint(URL)
	if err != nil {
		return nil, err
	}
	c, err := client.NewClient(ep)
	if err != nil {
		return nil, err
	}
	s, err := c.NewUploadPackSession(ep, nil)
	if err != nil {
		return nil, err
	}
	return s.AdvertisedReferences()
}

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
		u.Scheme = defaultScheme
	}

	options := &git.CloneOptions{}
	if u.Fragment != "" {
		options.ReferenceName = plumbing.ReferenceName("refs/heads/" + u.Fragment)
		u.Fragment = ""
	}
	options.URL = u.String()
	return options, err
}

func addDefaultURLScheme(URL string) string {
	u, err := url.Parse(URL)
	if err != nil {
		return URL
	}
	if u.Scheme == "" {
		u.Scheme = defaultScheme
	}
	return u.String()
}
