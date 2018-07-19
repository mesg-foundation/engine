package importer

import (
	"io/ioutil"
	"path/filepath"
)

func readDockerfile(path string) (data []byte, err error) {
	file := filepath.Join(path, "Dockerfile")
	data, err = ioutil.ReadFile(file)
	return
}
