package importer

import (
	"io/ioutil"
	"path/filepath"
)

func readDockerfile(path string) ([]byte, error) {
	file := filepath.Join(path, "Dockerfile")
	return ioutil.ReadFile(file)
}
