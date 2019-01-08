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
