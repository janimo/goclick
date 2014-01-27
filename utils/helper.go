//
// Copyright (c) 2014 Canonical Ltd.
//
// Author: Jani Monoses <jani@ubuntu.com>
//
package utils

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Alternative to filepath.Walk that does not ignore symlinks
type walkFunc func(string, os.FileInfo)

func Walk(path string, fn walkFunc, ignores []string) {
	for _, pattern := range ignores {
		ok, err := filepath.Match(pattern, filepath.Base(path))
		if err != nil {
			ExitError(err)
		}
		if ok {
			return
		}
	}
	fi, err := os.Lstat(path)
	if err != nil {
		ExitError(err)
	}

	fn(path, fi)

	if fi.IsDir() {
		fis, err := ioutil.ReadDir(path)
		if err != nil {
			ExitError(err)
		}
		for _, f := range fis {
			Walk(filepath.Join(path, f.Name()), fn, ignores)
		}
	}
}

// Exists tells whether a file named by a path exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// copyfile copies src to dest.
// if dest is be a regular file its content is copied from src
// if dest is a symlink or a directory it is created
func copyfile(src, dest string) {
	fi, err := os.Lstat(src)
	if err != nil {
		ExitError(err)
	}

	//Create an empty directory at dest
	if fi.IsDir() {
		err := os.Mkdir(dest, fi.Mode())
		if err != nil {
			ExitError(err)
		}
		return
	}

	//Create a symlink at dest
	if fi.Mode()&os.ModeSymlink != 0 {
		ln, err := os.Readlink(src)
		if err != nil {
			ExitError(err)
		}
		err = os.Symlink(ln, dest)
		if err != nil {
			ExitError(err)
		}
		return
	}

	d, err := os.Create(dest)
	if err != nil {
		ExitError(err)
	}
	defer d.Close()

	s, err := os.Open(src)
	if err != nil {
		ExitError(err)
	}
	defer s.Close()

	_, err = io.Copy(d, s)
	if err != nil {
		ExitError(err)
	}

	//Set the same file mode as src has
	err = d.Chmod(fi.Mode())
	if err != nil {
		ExitError(err)
	}
}

//cpa copies src to dest recursively, ignoring entries matching any of the ignore patterns
func Copytree(src, dest string, ignores []string) {
	dest = dest + "/"
	Walk(src, func(s string, fi os.FileInfo) {
		d := strings.Replace(s, src, dest, 1)
		copyfile(s, d)
	}, ignores)
}

// Get total size of all files in a dir
//Output of "du -ks --apparent-size"
func GetDirSize(dir string) string {
	size := 0
	filepath.Walk(dir, func(fn string, fi os.FileInfo, err error) error {
		size += int(fi.Size())
		return nil
	})
	return strconv.Itoa(size / 1024)
}

//Ensure the path ends in "/"
func AppendSlash(path string) string {
	if path[len(path)-1] != '/' {
		path = path + "/"
	}
	return path
}
