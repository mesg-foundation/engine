package xos

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Touch creates a new file, truncating it if it already exists.
func Touch(name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	return f.Close()
}

// Exist return true if given file exists, false otherwise.
func Exist(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// Copy copies file, symlink and file metadata from source to target path.
func Copy(src, dest string) error {
	si, err := os.Lstat(src)
	if err != nil {
		return err
	}

	if si.IsDir() {
		return CopyDir(src, dest)
	}

	// Handle symbolic link.
	if si.Mode()&os.ModeSymlink != 0 {
		return CopySymlink(src, dest)
	}

	if err := CopyFile(src, dest); err != nil {
		return err
	}

	// Set back file information.
	if err := os.Chtimes(dest, si.ModTime(), si.ModTime()); err != nil {
		return err
	}

	return os.Chmod(dest, si.Mode())
}

// CopyDir copies direcotry from source to target path.
func CopyDir(src, dest string) error {
	si, err := os.Lstat(src)
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dest, si.Mode()); err != nil {
		return err
	}

	for _, e := range entries {
		newsrc := filepath.Join(src, e.Name())
		newdst := filepath.Join(dest, e.Name())
		if err := Copy(newsrc, newdst); err != nil {
			return err
		}
	}
	return nil
}

// CopyFile copies file from source to target path.
func CopyFile(src, dest string) error {
	sr, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sr.Close()

	dw, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer dw.Close()

	_, err = io.Copy(dw, sr)
	return err
}

// CopySymlink copies file under symbolic link.
func CopySymlink(src, dest string) error {
	target, err := os.Readlink(src)
	if err != nil {
		return err
	}
	return os.Symlink(target, dest)
}
