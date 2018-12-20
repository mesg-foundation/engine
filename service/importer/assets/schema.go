// Code generated by go-bindata. (@generated) DO NOT EDIT.
// sources:
// service/importer/assets/schema.json
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

var _serviceImporterAssetsSchemaJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x57\x5d\x4f\xdb\x3c\x14\xbe\xcf\xaf\xb0\x0c\x77\xef\x5b\xca\x34\x69\x12\xbd\xdb\xee\x26\x4d\x03\x89\xbb\x41\x36\x99\xe4\xa4\x98\x25\xb6\xf1\x07\x5a\x87\xfa\xdf\xa7\x38\x6d\xea\xc4\x76\x12\xa0\x43\xd3\xba\x5e\xb5\x3e\x3e\x1f\xcf\x73\x8e\x1f\xbb\x8f\x09\x42\xf8\x58\x65\xb7\x50\x11\xbc\x40\xf8\x56\x6b\xb1\x98\xcf\xef\x14\x67\xb3\x66\xf5\x84\xcb\xe5\x3c\x97\xa4\xd0\xb3\xd3\x77\xf3\x66\xed\x08\xff\x5f\xfb\xe9\x95\x80\xda\x89\xdf\xdc\x41\xa6\x9b\x35\x09\xf7\x86\x4a\xc8\xf1\x02\x5d\x25\x08\x21\x84\x19\xa9\x00\x27\x08\xa5\xd6\x4e\xf2\x9c\x6a\xca\x19\x29\x2f\x24\x17\x20\x35\x05\x85\x17\xa8\x20\xa5\x02\xbb\x41\xb8\xcb\x8f\x4e\x88\xed\x2f\x27\xb1\xd2\x92\xb2\xa5\x4d\x6c\xd7\x2b\xca\x3e\x01\x5b\xea\x5b\xbc\x40\x6f\xec\xe2\xba\xb1\x61\x45\xf3\x49\x01\xc8\x8f\x36\xc0\xdb\xb3\x76\x59\x10\xad\x41\xb2\xda\xe3\xeb\xd5\xe9\xec\x8c\xcc\x7e\xbe\x9f\x7d\xb9\xbe\x3e\xf9\x36\x4b\xff\x3b\xc6\x9d\x4c\x39\xa8\x4c\x52\x51\x63\x1c\xc8\xd8\x71\x91\x20\xb8\xa2\x9a\xcb\xd5\x54\x0f\x78\x00\xa6\x95\xbb\x9b\x33\x38\x2f\x5a\xd2\xeb\xcf\xe3\x2e\x04\x33\x65\x89\xb7\xce\xd6\xd6\x7e\x0b\xb7\xb1\xb5\x8d\xb5\xab\xdd\xb8\xa1\xe8\xc2\xef\x5e\xbb\xc5\xa1\xae\xe1\xad\xbf\xa3\x1e\x45\x09\x35\x0a\x7c\x34\xcf\xa1\xa0\xcc\xe6\x56\x73\x0b\x17\x77\xf6\xae\x93\xd0\xf7\xed\xb7\xb4\xc3\x96\x26\xea\xfb\xe1\x90\x55\xa3\x7d\x3e\x57\x39\x08\x60\x39\xb0\xac\x5b\x55\x0c\xf7\x24\xcc\x23\x78\xc7\xb0\x46\x70\xb6\x95\xae\xba\x5d\xf0\x14\xa8\xb5\xd0\x8a\x2c\xc1\xa5\x26\xf5\xc8\x58\x77\xc8\xc8\x38\x2b\xe8\xd2\x48\xd2\x3f\xcd\xa3\x15\x25\xdb\x60\x36\x14\x76\x76\xed\x54\x4d\x10\x49\x2a\xd0\x20\xff\xc6\xd9\x1c\xa8\xe5\x49\xf5\x0c\xb7\x74\x97\xaa\xb7\x9c\x7a\x31\x02\xf7\x4a\xc7\xde\xbb\x63\x82\x60\x5c\x29\x76\x3f\xeb\x7e\xb6\xe8\x3d\xf0\xe2\xb0\x5c\x34\xac\x8d\xc5\xbc\xe1\xbc\x04\xc2\xa6\x05\xdd\x38\x4d\x29\xd2\xf7\xae\xaf\x23\x66\xaa\x60\x6f\xac\xf5\x32\xee\x89\x10\xfe\x6c\xaa\x1b\x90\x31\xeb\x87\x0d\x8c\x88\xf9\xbc\x99\xae\x80\xd1\x1b\x00\xd4\x34\xa5\x20\xa6\xd4\x35\x9e\xcb\x18\xe9\xc9\xd0\xef\xa7\x5e\x3b\x7b\x93\xd0\xe0\x11\xc0\xdc\x68\x61\xb4\xda\xa2\x48\x77\x8a\x1b\x91\xda\xc0\x94\x47\xc7\xd0\x99\x94\x81\x69\x9e\xe4\x4e\x99\x2d\x73\x92\xae\x3b\xb2\x18\x8c\xa5\x20\x93\xb0\xaf\x60\x5b\x02\x23\x98\x0e\x4c\x49\x73\xa2\xc9\xc1\x28\xa9\x05\x1b\x89\x37\x7d\x98\xda\x0c\x2f\x90\x8d\xb5\xff\xb6\xff\xcd\xba\xe1\x36\xfa\x8f\x15\x8d\x40\x83\x9e\x72\xca\x83\xe4\x3a\x4f\xb4\xbd\x3d\x6e\x23\xac\x35\xaf\xcd\x61\xdc\x1d\xbd\xf0\xff\xbd\xf6\x19\xc9\x78\x55\x11\x96\x3f\x87\x4c\x22\x97\x9e\xca\xf9\x4f\x4d\x34\xf2\xdc\x44\xbd\x27\x67\x27\x3d\x91\x92\xac\x7c\x7d\xa2\x1a\xaa\x88\x52\x0c\x1f\xdb\xf8\x11\x4a\x83\x10\x1f\x78\x69\x2a\x4f\x94\x5e\x09\xa5\x61\xf4\xde\xc0\xc7\x0d\x56\x2d\x8d\xaf\xbd\xaf\x4d\x44\x21\x79\xf5\x8f\x0c\xc1\xa5\x7f\xbf\x1f\x1e\x0d\xc0\x1e\x62\xb2\xe1\x57\x3f\x52\x79\xa4\xea\x81\x8a\xa3\x77\x5e\xfd\x0f\x39\x59\x27\xbf\x02\x00\x00\xff\xff\x19\x00\xfb\x06\x87\x14\x00\x00")

func serviceImporterAssetsSchemaJsonBytes() ([]byte, error) {
	return bindataRead(
		_serviceImporterAssetsSchemaJson,
		"service/importer/assets/schema.json",
	)
}

func serviceImporterAssetsSchemaJson() (*asset, error) {
	bytes, err := serviceImporterAssetsSchemaJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "service/importer/assets/schema.json", size: 5255, mode: os.FileMode(420), modTime: time.Unix(1544790675, 0)}
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
	"service/importer/assets/schema.json": serviceImporterAssetsSchemaJson,
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
		"importer": &bintree{nil, map[string]*bintree{
			"assets": &bintree{nil, map[string]*bintree{
				"schema.json": &bintree{serviceImporterAssetsSchemaJson, map[string]*bintree{}},
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
