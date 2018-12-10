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

var _serviceImporterAssetsSchemaJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x57\x4f\x4f\xdc\x3e\x10\xbd\xe7\x53\x58\x86\xdb\xef\xb7\x2c\x55\xa5\x4a\xec\xad\xbd\x55\xaa\x0a\x12\xb7\x42\x5a\x99\x64\xb2\x98\x26\xb6\xf1\x1f\xd4\x2d\xda\xef\x5e\xc5\xd9\xcd\x3a\xb1\x9d\x04\xd8\xa2\xaa\xdb\x9c\x12\x8f\x67\xc6\xef\xcd\xf8\xd9\x79\x4c\x10\xc2\xc7\x2a\xbb\x85\x8a\xe0\x05\xc2\xb7\x5a\x8b\xc5\x7c\x7e\xa7\x38\x9b\x35\xa3\x27\x5c\x2e\xe7\xb9\x24\x85\x9e\x9d\xbe\x9b\x37\x63\x47\xf8\xff\xda\x4f\xaf\x04\xd4\x4e\xfc\xe6\x0e\x32\xdd\x8c\x49\xb8\x37\x54\x42\x8e\x17\xe8\x2a\x41\x08\x21\xcc\x48\x05\x38\x41\x28\xb5\x76\x92\xe7\x54\x53\xce\x48\x79\x21\xb9\x00\xa9\x29\x28\xbc\x40\x05\x29\x15\xd8\x09\xc2\x1d\x7e\x74\x42\x6c\xbf\x9c\xc4\x4a\x4b\xca\x96\x36\xb1\x1d\xaf\x28\xfb\x04\x6c\xa9\x6f\xf1\x02\xbd\xb1\x83\xeb\xc6\x86\x15\xcd\x27\x05\x20\x3f\xda\x00\x6f\xcf\xda\x61\x41\xb4\x06\xc9\x6a\x8f\xaf\x57\xa7\xb3\x33\x32\xfb\xf9\x7e\xf6\xe5\xfa\xfa\xe4\xdb\x2c\xfd\xef\x18\x77\x32\xe5\xa0\x32\x49\x45\x8d\x71\x20\x63\xc7\x45\x82\xe0\x8a\x6a\x2e\x57\x53\x3d\xe0\x01\x98\x56\xee\x6c\xce\xe0\xbc\x68\x49\xaf\x9f\xc7\x5d\x08\x66\xca\x12\x6f\x9d\xad\xad\x7d\x0b\x97\xb1\xb5\x8d\x95\xab\x9d\xb8\xa1\xe8\xc2\xaf\x5e\x3b\xc5\xa1\xae\xe1\xad\x3f\xa3\x6e\x45\x09\x35\x0a\x7c\x34\xcf\xa1\xa0\xcc\xe6\x56\x73\x0b\x17\x77\xe6\xae\x93\xd0\xfb\xf6\x2d\xed\xb0\xa5\x89\xfa\x7e\x38\x64\xd5\x68\x9f\xcf\x55\x0e\x02\x58\x0e\x2c\xeb\xae\x2a\x86\x7b\x12\xe6\x11\xbc\x63\x58\x23\x38\xdb\x95\xae\xba\x55\xf0\x14\xa8\xb5\xd0\x8a\x2c\xc1\xa5\x26\xf5\xc8\x58\x77\xc8\xc8\x38\x2b\xe8\xd2\x48\xd2\xdf\xcd\xa3\x2b\x4a\xb6\xc1\x6c\x28\xec\xcc\xda\xa9\x9a\x20\x92\x54\xa0\x41\xfe\x8d\xbd\x39\xb0\x96\x27\xad\x67\xb8\xa4\xbb\x54\xbd\xe1\xd4\x8b\x11\x38\x57\x3a\xf6\xde\x19\x13\x04\xe3\x4a\xb1\xfb\xac\xfb\xd9\xa2\xe7\xc0\x8b\xc3\x72\xd1\xb0\x36\x16\xf3\x86\xf3\x12\x08\x9b\x16\x74\xe3\x34\x65\x91\xbe\x77\x7d\x1c\x31\x53\x05\x6b\x63\xad\x97\x71\x4f\x84\xf0\x67\x53\xdd\x80\x8c\x59\x3f\x6c\x60\x44\xcc\xe7\x4d\x77\x05\x8c\x5e\x03\xa0\xa6\x28\x05\x31\xa5\xae\xf1\x5c\xc6\x48\x4f\x86\xbe\x9f\x7a\xec\xec\x4d\x42\x83\x5b\x00\x73\xa3\x85\xd1\x6a\x8b\x22\xdd\x29\x6e\x44\x6a\x03\x5d\x1e\x6d\x43\xa7\x53\x06\xba\x79\x92\x3b\x65\x76\x99\x93\x74\xdd\x91\xc5\x60\x2c\x05\x99\x84\x7d\x05\xdb\x12\x18\xc1\x74\x60\x4a\x9a\x13\x4d\x0e\x46\x49\x2d\xd8\x48\xbc\xe9\xcd\xd4\x66\x78\x81\x6c\xac\xfd\xbb\xfd\x6f\xd6\x0d\xb7\xd0\x7f\xac\x68\x04\x0a\xf4\x94\x5d\x1e\x24\xd7\xb9\xa2\xed\xed\x72\x1b\x61\xad\xb9\x6d\x0e\xe3\xee\xe8\x85\xff\xf7\xda\x67\x24\xe3\x55\x45\x58\xfe\x1c\x32\x89\x5c\x7a\x2a\xe7\x5f\x35\xd1\xc8\x75\x13\xf5\xae\x9c\x9d\xf4\x44\x4a\xb2\xf2\xf5\x89\x6a\xa8\x22\x4a\x31\xbc\x6d\xe3\x5b\x28\x0d\x42\x7c\xe0\xa5\xa9\x3c\x51\x7a\x25\x94\x86\xd1\x7b\x03\x1f\x37\x58\xb5\x34\xbe\xf6\xbe\x36\x11\x85\xe4\xd5\x3f\x32\x04\x97\xfe\xf9\x7e\x20\x34\x74\x25\xb0\xfe\x1b\x4d\xd6\xc9\xaf\x00\x00\x00\xff\xff\xa2\xc0\xe3\xdf\xf3\x13\x00\x00")

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

	info := bindataFileInfo{name: "service/importer/assets/schema.json", size: 5107, mode: os.FileMode(420), modTime: time.Unix(1544151489, 0)}
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
