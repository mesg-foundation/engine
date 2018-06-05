// Code generated by go-bindata.
// sources:
// service/schema.json
// DO NOT EDIT!

package service

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

var _serviceSchemaJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x58\x5f\x6f\xea\x36\x14\x7f\xe7\x53\x58\x2e\x0f\xed\x5a\x4a\xf7\x32\x69\xbc\x54\x5d\x85\xa6\x6a\x1d\x54\xb0\x6a\xd2\x10\x9d\x4c\x72\x02\xee\x12\x3b\x75\x1c\x34\xd6\xf1\xdd\xaf\x92\x90\x60\x27\x36\x09\x6d\xa9\xae\x2e\x17\xe9\x4a\xb9\xe7\xff\x39\x3e\xfe\x9d\xe3\xbe\xb6\x10\xc2\xed\xc8\x59\x40\x40\x70\x0f\xe1\x85\x94\x61\xaf\xdb\x7d\x8e\x38\xeb\x64\xd4\x4b\x2e\xe6\x5d\x57\x10\x4f\x76\xae\x7e\xea\x66\xb4\x13\x7c\x91\xe8\xc9\x55\x08\x89\x12\x9f\x3d\x83\x23\x33\x9a\x80\x97\x98\x0a\x70\x71\x0f\x4d\x5a\x08\x21\x84\x19\x09\x00\xb7\x10\x9a\xa6\x7c\xe2\xba\x54\x52\xce\x88\xff\x20\x78\x08\x42\x52\x88\x70\x0f\x79\xc4\x8f\x20\x15\x08\x55\xf2\x6b\x66\x62\x09\x22\xa2\x9c\x15\x04\xc5\x77\x24\x05\x65\x73\x9c\x92\xd7\x17\x8a\x47\xbb\xec\x45\x4e\x0f\x28\xbb\x07\x36\x97\x0b\xdc\x43\x3f\x6a\x16\x5c\x88\x1c\x41\x43\xb9\x87\xd3\x25\x8d\xe8\x8c\xfa\x54\xae\x54\x8d\xb6\x00\x2f\xd1\x38\xe9\xba\xe0\x51\x96\xa6\x1e\x75\x15\x59\xcd\x46\x18\xcf\x7c\x1a\x2d\xea\x0d\xe4\x82\x9a\x36\x2c\x81\xc9\x48\x55\xe6\x0c\x86\x5e\x71\x12\xc9\xef\x75\x9b\x04\x8b\x7d\x1f\xe7\xca\x29\xaf\xf8\x32\x9f\x6d\xc1\xab\x3b\xc3\x42\x30\x24\x52\x82\x60\x0f\xd5\x23\x2d\x44\x9e\x26\xa4\xf3\xdf\x4d\xe7\xaf\xbf\x3b\xd3\xf3\x76\x85\x6d\x2d\x40\x9a\x2b\xd6\x64\xd7\x2d\xd3\x77\xfe\x35\xd5\x4a\x25\x49\xf4\xcf\x91\x54\x2a\x49\xf5\xed\x85\x72\x21\x04\xe6\x02\x73\xf4\x90\x6c\x49\x37\x4a\xb8\x26\xd9\x9d\x89\x5a\x92\x2c\xc2\x5c\xe9\xf5\xaf\xa0\x51\xc1\xa1\x01\x99\x83\x5a\x97\x69\xa5\x12\x6b\xad\x12\x0e\x67\x1e\x9d\xc7\x82\x94\x41\xa1\x36\xa2\x56\x6e\x2c\x35\x85\x15\x29\x05\xe1\x8c\xe8\x61\x03\x2e\x17\x3c\x12\xfb\x32\x61\xdd\xdc\xdf\x6f\xe9\xc0\xe2\x40\x4b\x15\x2b\xe5\xd0\x44\x11\xc2\x8f\xe3\xfe\x68\xac\x12\xfe\x1c\x8e\x7e\x2b\x91\x06\xc3\x41\x1f\x1b\x1b\xc3\x00\x55\x07\x8d\x76\x3c\x7c\x1c\xdd\xf6\x55\xca\xed\x70\xf0\xc7\xcd\xdd\xa0\x3f\x6a\x18\x30\x11\x24\x00\x09\xe2\x9b\xbb\xf7\x3b\x02\xd9\x2b\x18\xb4\xf3\xc6\x6c\x5d\x95\xc8\xd3\x8a\x0d\xc3\x08\xd7\xf8\xa5\xf9\x6c\x4c\x46\x1d\xb1\xea\x6f\x5d\xf6\x66\x9d\xd6\xef\x36\xcb\xc3\xac\x6a\x75\x36\x67\x9c\xfb\x40\x58\x33\xa3\x1b\xa5\x26\x41\x56\xb5\x4d\x97\x46\xe3\x8e\xed\x9a\xc9\xe5\x88\x83\x19\x08\x1b\xf7\x97\x4d\x1a\x16\xf6\x30\xeb\x2e\x03\xb3\xd2\x00\xa8\x74\xeb\xc7\xb6\xa2\xb7\x76\xfd\x7f\x9f\x31\xe5\xc1\x47\x8e\x27\xcb\x5c\x82\x80\xca\x12\x7e\xa0\x9a\x63\xcb\xef\x77\xc2\x7e\x6a\xff\xff\x74\x3a\xb9\xea\xfc\x3c\xfd\x61\x72\x39\xbd\x4e\xbf\xce\xcf\x4e\xd1\xe9\xef\xfd\xf1\xaf\x67\x67\xd7\xed\x6d\x81\x94\xc6\xc1\x2e\x2c\xc1\x4f\x02\xfa\x5c\xb7\xf0\x2f\x38\xb1\xe4\x9f\xec\x75\x49\x7c\xea\x92\xcf\x71\x6b\x9c\xf2\xe9\xb6\xf4\x51\x8d\x64\xc4\x52\xcc\x63\x19\xc6\x32\x2a\xa6\x54\x6d\xdf\x19\xe0\xd2\x8a\x67\x7a\xe7\xd8\x60\xb1\x91\xfa\x12\x04\xf5\x28\x99\xf9\x56\xe7\x15\xe0\x53\xd5\x43\xb2\x7a\xb3\x6e\xe9\x3a\x23\xfb\x9a\x95\x4a\x36\x5c\xf9\x8a\x96\x46\x08\x99\x17\x3f\x25\x04\xca\xd2\x43\x6a\x14\x84\xb2\x5a\x18\x6d\x45\xe0\x08\xf8\x28\x63\x79\xfb\x58\xea\x7a\x4c\xdb\x88\x4b\x24\x39\x9a\x6d\x24\x4d\xd6\x62\xaf\x79\x27\x15\x1e\xde\x31\x7a\x75\xc4\xcc\x5e\xe2\x87\x85\x4c\xf5\xa0\xbf\x5a\xbc\x3c\x10\x68\xe5\x4b\x47\x13\xd0\x32\xf4\xc8\x3e\x28\x63\x3c\x5f\xe5\x21\x7b\xe8\x05\x2b\x7b\x93\xef\x31\xf9\xab\x7f\xc0\x2b\x57\xc4\xe1\x41\x40\x98\xfb\xa6\xf9\xc7\xfd\x38\xa8\x1e\x69\xf5\xc5\x88\x6a\x5e\x8d\xa8\xf4\x72\xd4\x22\x20\x42\x90\x55\x15\x25\x63\x46\x5f\x62\xb8\x93\x10\x24\x01\x48\x11\x57\xb1\x91\x6e\x98\x06\x48\xdb\x8d\x2f\xf6\xbb\x6e\xee\xaa\x4d\x21\x3c\xc1\x83\xef\xc5\x08\xb9\xa8\x0e\xdf\x23\x29\x83\x8e\x0f\xad\xe4\xdf\xfa\x4b\x00\x00\x00\xff\xff\x55\xd7\xe8\x6f\x41\x18\x00\x00")

func serviceSchemaJsonBytes() ([]byte, error) {
	return bindataRead(
		_serviceSchemaJson,
		"service/schema.json",
	)
}

func serviceSchemaJson() (*asset, error) {
	bytes, err := serviceSchemaJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "service/schema.json", size: 6209, mode: os.FileMode(420), modTime: time.Unix(1527897332, 0)}
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
	"service/schema.json": serviceSchemaJson,
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
	"service": &bintree{nil, map[string]*bintree{
		"schema.json": &bintree{serviceSchemaJson, map[string]*bintree{}},
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

