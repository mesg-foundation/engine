package readme

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
)

var readmeRegexp = regexp.MustCompile(`(?i)^readme(\\.md)?`)

// Lookup returns the content of a readme in a directory
func Lookup(path string) (string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if readmeRegexp.Match([]byte(file.Name())) {
			readme, err := ioutil.ReadFile(filepath.Join(path, file.Name()))
			if err == nil && len(readme) > 0 {
				return string(readme), nil
			}
		}
	}
	return "", nil
}
