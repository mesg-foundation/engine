package xos

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestTouch(t *testing.T) {
	name := filepath.Join(os.TempDir(), "__test-file")
	defer Remove(name)

	if err := Touch(name); err != nil {
		t.Fatalf("Touch failed: %s", err)
	}
	if !Exist(name) {
		t.Fatalf("Touched file dosen't exist")
	}

	if err := Touch(name); err != nil {
		t.Fatalf("Touch failed: %s", err)
	}
}

func TestExist(t *testing.T) {
	name := filepath.Join(os.TempDir(), "__test-file")
	defer Remove(name)
	if err := Touch(name); err != nil {
		t.Fatalf("Touch failed: %s", err)
	}

	if !Exist(name) {
		t.Fatalf("Exist got: %t, want: %t ", false, true)
	}
	if Exist(name + "0") {
		t.Fatalf("Exist got: %t, want: %t ", true, false)
	}
}

func TestCopyDir(t *testing.T) {
	var (
		srcdir    = filepath.Join(os.TempDir(), "__test-dir-src")
		dstdir    = filepath.Join(os.TempDir(), "__test-dir-dst")
		srcsubdir = filepath.Join(srcdir, "a", "b")
		file      = filepath.Join(srcsubdir, "__test-file")
		content   = "test"
	)
	defer Remove(srcdir, dstdir, file)

	if err := os.MkdirAll(srcsubdir, 0777); err != nil {
		t.Fatalf("Mkdir failed: %s", err)
	}

	if err := ioutil.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile failed: %s", err)
	}

	if err := CopyDir(srcdir, dstdir); err != nil {
		t.Fatalf("CopyDir failed: %s", err)
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatalf("ReadFile failed: %s", err)
	}

	if string(b) != content {
		t.Fatalf("Copy failed with different content - got: %s, want: %s", string(b), content)
	}
}

func TestCopyFile(t *testing.T) {
	content := "test"
	srcname := filepath.Join(os.TempDir(), "__test-file_src")
	dstname := filepath.Join(os.TempDir(), "__test-file_dst")
	defer Remove(srcname, dstname)

	if err := ioutil.WriteFile(srcname, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile failed: %s", err)
	}

	if err := CopyFile(srcname, dstname); err != nil {
		t.Fatalf("Copy failed: %s", err)
	}

	b, err := ioutil.ReadFile(dstname)
	if err != nil {
		t.Fatalf("ReadFile failed: %s", err)
	}

	if string(b) != content {
		t.Fatalf("Copy failed with different content - got: %s, want: %s", string(b), content)
	}
}

func TestCopySymlink(t *testing.T) {
	name := filepath.Join(os.TempDir(), "__test-file_src")
	srcsymlink := filepath.Join(os.TempDir(), "__test-symlink_src")
	dstsymlink := filepath.Join(os.TempDir(), "__test-symlink_dst")
	defer Remove(srcsymlink, name, dstsymlink)

	if err := Touch(name); err != nil {
		t.Fatalf("touch file failed: %s", err)
	}

	if err := os.Symlink(name, srcsymlink); err != nil {
		t.Fatalf("symlink to file failed: %s", err)
	}

	if err := CopySymlink(srcsymlink, dstsymlink); err != nil {
		t.Fatalf("CopySymlink failed: %s", err)
	}

	fi, err := os.Lstat(dstsymlink)
	if err != nil {
		t.Fatalf("file info failed: %s", err)
	}

	if fi.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("new file is not symbolic link")
	}
}
