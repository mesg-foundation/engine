package readme

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
)

// LookupReadme returns the content of a readme in a directory
func LookupReadme(path string) (string, error) {
	reg := regexp.MustCompile("(?i)^readme(\\.md)?")
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if reg.Match([]byte(file.Name())) {
			readme, err := ioutil.ReadFile(filepath.Join(path, file.Name()))
			if err == nil && len(readme) > 0 {
				return string(readme), nil
			}
		}
	}
	return "", nil
}
