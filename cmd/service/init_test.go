package service

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

func TestWriteInFolder(t *testing.T) {
	dir, _ := ioutil.TempDir("", "TestWriteInFolder")
	os.RemoveAll(dir)
	err := writeInFolder(dir, []byte("hello world"))
	defer os.RemoveAll(dir)
	assert.Nil(t, err)
	mesg, err := ioutil.ReadFile(filepath.Join(dir, "mesg.yml"))
	assert.Nil(t, err)
	assert.Equal(t, string(mesg), "hello world")
	_, err = ioutil.ReadFile(filepath.Join(dir, "Dockerfile"))
	assert.Nil(t, err)
}

func TestWriteInFolderExistingFolder(t *testing.T) {
	dir, _ := ioutil.TempDir("", "TestWriteInFolder")
	err := writeInFolder(dir, []byte("hello world"))
	defer os.RemoveAll(dir)
	assert.NotNil(t, err)
}

func TestGenerateMesgFile(t *testing.T) {
	file, err := generateMesgFile(&service.Service{
		Name:        "TestGenerateMesgFile",
		Description: "description",
	})
	assert.Nil(t, err)
	res := strings.Replace(strings.Replace(templateText, "{{.Name}}", "TestGenerateMesgFile", -1), "{{.Description}}", "description", -1)
	assert.Equal(t, string(file), res)
}
