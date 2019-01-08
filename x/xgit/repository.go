// Copyright 2018 MESG Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
