// Code generated by go-bindata.
// sources:
// cmd/service/assets/readmeTemplate.md
// DO NOT EDIT!

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

var _cmdServiceAssetsReadmetemplateMd = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x92\xbf\x4e\xc3\x30\x10\xc6\x77\x3f\xc5\x49\xe9\x48\xf2\x00\x48\x0c\xfc\xe9\x50\x81\x08\x2a\x65\xe8\x94\x98\xe6\x28\x51\x69\x12\xc5\x6e\xa5\xc8\xf1\xbb\xa3\xb3\x9d\xba\xb6\x0a\x4b\x97\xe4\x3e\xfb\x3b\xdf\xe7\x9f\x9c\x80\x52\xd9\x2b\xdf\xa3\xd6\x8c\x29\x95\x3d\xa1\xd8\xf4\x75\x27\xeb\xb6\xa1\x95\xb2\x2c\x3f\xb9\xf8\x66\x7b\x14\xdb\x74\xd3\xf6\x08\x02\xfb\x63\xbd\x41\xa8\xb0\xfb\x69\x07\x28\x8a\xe5\xfc\xed\xe5\xfe\x71\x5e\x3c\xac\x8b\x75\xfe\xb1\x24\x9d\xbf\x2f\x56\xf9\x72\x5d\x14\xd4\x4f\xc7\xd6\x5f\x90\xcd\x8f\xd8\x48\xa1\x35\x4b\x12\xb0\x35\x53\xaa\xe7\xcd\x16\x61\xb6\xc3\xe1\x06\x66\x48\xab\x70\x7b\x17\x78\x29\x20\xed\xdb\x7c\xd6\x13\xa7\x34\x03\xa6\x2d\x2e\xb9\xd6\x23\x34\x7c\x8f\x30\x82\x1c\x3a\xfa\x55\xbe\x01\x46\x36\x42\x9a\xa6\x10\x7c\x7d\x96\x8a\x4b\xfe\x6c\xf2\x50\x45\x71\xa2\x93\x95\x9a\x3c\x5a\xc3\x08\xa5\xd3\xd9\x6a\xe8\x50\xeb\x12\x4e\x8e\x30\xa6\x99\x81\x4d\xa5\xf5\xc5\x9f\xa3\xb4\xe2\x62\xe7\x20\x99\x32\x66\x24\xb9\xd8\x19\x44\xde\x18\x11\x22\x47\x0c\xc8\x11\x32\x5b\x8b\xa6\x3b\x10\xdb\x84\x3a\xad\x60\xec\x0a\x5e\x35\x1d\x61\x81\x99\xd2\x10\x0b\x26\x19\x20\x93\xcd\x30\x9b\xb4\x63\x76\xbe\xf2\x3f\xb3\xe0\x26\xf9\x41\x9e\x5d\xc5\x29\xe6\x93\xb5\x66\xc5\x46\xb3\xb5\xcf\x16\xf4\x1a\x82\x27\xb7\xe3\x68\x75\x4c\xf2\x0a\x50\x7f\xc4\x99\xe6\xf8\xd7\x75\x96\xc4\xbd\x2f\xe7\x09\x5e\xd8\xa5\x7c\x9e\xd7\x05\x6e\x81\xfe\x0d\x00\x00\xff\xff\xa8\x37\xde\x3c\xf8\x03\x00\x00")

func cmdServiceAssetsReadmetemplateMdBytes() ([]byte, error) {
	return bindataRead(
		_cmdServiceAssetsReadmetemplateMd,
		"cmd/service/assets/readmeTemplate.md",
	)
}

func cmdServiceAssetsReadmetemplateMd() (*asset, error) {
	bytes, err := cmdServiceAssetsReadmetemplateMdBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "cmd/service/assets/readmeTemplate.md", size: 1016, mode: os.FileMode(420), modTime: time.Unix(1530260683, 0)}
	a := &asset{bytes: bytes, info: info}
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
	"cmd": &bintree{nil, map[string]*bintree{
		"service": &bintree{nil, map[string]*bintree{
			"assets": &bintree{nil, map[string]*bintree{
				"readmeTemplate.md": &bintree{cmdServiceAssetsReadmetemplateMd, map[string]*bintree{}},
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

