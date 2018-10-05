package xos

import "os"

// Remove removes all given named file or directory.
func Remove(names ...string) error {
	var err error
	for _, name := range names {
		if err1 := os.Remove(name); err == nil {
			err = err1
		}
	}
	return err
}

// RemoveAll removes all given path and any children it contains.
func RemoveAll(paths ...string) error {
	var err error
	for _, path := range paths {
		if err1 := os.RemoveAll(path); err == nil {
			err = err1
		}
	}
	return err
}
