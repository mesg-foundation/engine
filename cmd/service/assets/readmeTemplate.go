// Code generated by go-bindata. DO NOT EDIT.
// sources:
// cmd/service/assets/readmeTemplate.md
package assets

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _cmdServiceAssetsReadmetemplateMd = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x92\xc1\x6e\x83\x30\x0c\x86\xef\x79\x0a\x4b\xf4\x84\x06\x0f\x50\x69\xb7\xed\x30\x55\xda\xa4\xa9\x0f\x40\xd6\x7a\x1d\x62\x05\x44\x68\x25\x94\xfa\xdd\x27\x3b\x40\x12\xc6\x76\xea\xa5\x75\x7e\xdb\xf9\xff\x7c\x6d\x02\xd6\xe6\xaf\xfa\x8c\x44\x4a\x59\x9b\x3f\xa1\x39\x74\x65\xdb\x97\x4d\x4d\xa4\xac\x2d\x3f\x21\x7f\xc7\xb6\x31\x65\xdf\x74\x03\x91\x2a\x8a\xe2\x43\x9b\x2f\x75\x46\x73\xca\x0e\x4d\x87\x60\xb0\xbb\x96\x07\x84\x23\xb6\xdf\xcd\xc0\xf7\x2d\x17\x94\xb5\x58\x1f\xe7\xfb\x9e\xaf\x58\xf7\x86\x28\x01\x57\x29\x6b\x3b\x5d\x9f\x10\x36\x15\x0e\x0f\xb0\x41\x56\x61\xfb\xe8\x27\x55\xc2\x39\xb9\xcd\x31\x45\x85\x0a\x87\x2d\x14\x93\x5a\x70\x7a\xb7\xb9\x78\x83\x33\x9d\x5a\xba\xd7\x44\x37\x48\xd3\x1d\x0e\x69\x0a\x5c\xed\x87\x16\xc7\x32\xd8\x64\x45\xdd\x20\xcb\x32\x88\x3e\x7d\xd8\xa3\xee\xf5\x4e\x02\x73\xc5\x79\x7f\x99\x58\x3b\x4d\x11\x89\x45\x31\x2a\x39\x9b\x12\x15\x70\x83\x49\x89\x52\x8b\x8f\x30\x5b\xfd\x1a\x41\xee\xb5\xa9\x84\x0e\x48\xb5\xe4\xd8\x6b\x53\x09\xc6\x79\x2e\xa4\xc8\xe2\x2a\x44\x5e\x5b\x67\x28\x9d\x97\xba\xbd\xc8\xaf\x97\x24\xe0\x6a\xa5\xee\x03\xb4\xe4\xdb\x1c\x51\x29\x05\x69\xe4\x39\x32\x9d\x06\x03\xa8\x22\xc5\x54\x9d\xf4\x3f\xd6\xe8\x65\x6f\x97\xde\x3f\x6d\x3c\x28\x1f\xaf\x11\xc5\xe5\x73\xb5\x0f\x18\xae\x3a\xca\xf3\x34\x9b\xb8\xb6\xa7\x1d\x34\x1d\x73\x27\x2c\xa9\xdf\x07\xeb\x1f\xb9\x27\xcb\xf0\xcf\x1a\xe4\x9a\xc9\x8e\x73\x11\xda\xb5\xb8\x9e\xed\x0a\xe3\xe8\xfc\x13\x00\x00\xff\xff\x03\xb4\x0d\x6e\x74\x04\x00\x00")

func cmdServiceAssetsReadmetemplateMdBytes() ([]byte, error) {
	return bindataRead(
		_cmdServiceAssetsReadmetemplateMd,
		"cmd/service/assets/readmeTemplate.md",
	)
}

func cmdServiceAssetsReadmetemplateMd() (*asset, error) {
	readmeBytes, err := cmdServiceAssetsReadmetemplateMdBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "cmd/service/assets/readmeTemplate.md", size: 1140, mode: os.FileMode(420), modTime: time.Unix(1531285359, 0)}
	a := &asset{bytes: readmeBytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"cmd/service/assets/readmeTemplate.md": cmdServiceAssetsReadmetemplateMd,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"cmd": {nil, map[string]*bintree{
		"service": {nil, map[string]*bintree{
			"assets": {nil, map[string]*bintree{
				"readmeTemplate.md": {cmdServiceAssetsReadmetemplateMd, map[string]*bintree{}},
			}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

