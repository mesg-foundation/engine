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
