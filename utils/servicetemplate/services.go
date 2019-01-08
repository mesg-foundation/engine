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

package servicetemplate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mesg-foundation/core/x/xgit"
	"github.com/mesg-foundation/core/x/xos"
)

const templatesURL = "https://raw.githubusercontent.com/mesg-foundation/awesome/master/templates.json"

// Template represents single entry on awesome github list.
type Template struct {
	Name string
	URL  string
}

func (s *Template) String() string {
	return fmt.Sprintf("%s (%s)", s.Name, s.URL)
}

// List returns all service templates from awesome github project.
func List() ([]*Template, error) {
	resp, err := http.Get(templatesURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var list []*Template
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// Download download given template into dst directory.
func Download(t *Template, dst string) error {
	path, err := ioutil.TempDir("", "mesg-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(path)

	if err = xgit.Clone(t.URL, path); err != nil {
		return err
	}

	return xos.CopyDir(filepath.Join(path, "template"), dst)
}
